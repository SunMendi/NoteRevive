package auth

import (
	"fmt"
	"os"
	"strings"
	"time"

	"notemind/internal/llm"

	mailjet "github.com/mailjet/mailjet-apiv3-go"
	"gorm.io/gorm"
)

type AuthRepo interface {
	Create(user *User) error 
	GetByEmail(email string) (*User, error)
	SendDailySummary() error 
}

type authRepo struct {
	db  *gorm.DB
	llm *llm.LLMService
}

func NewAuthRepo(db *gorm.DB, llm *llm.LLMService) AuthRepo {
	return &authRepo{
		db:  db,
		llm: llm,
	}
}


func(r *authRepo) Create(user *User) error {
	 return r.db.Create(user).Error
}

func (r *authRepo) GetByEmail( email string ) (*User, error) {
	 var user User 
	 err := r.db.Where("email = ?", email).First(&user).Error 
	 if err != nil {
		 return nil, err 
	 }
	 return &user , nil 
}

func (r *authRepo) SendDailySummary() error{
	var users []User
	var finalMessage string
	var emailErrors []string
	var sentCount int

	err := r.db.Where("timezone IS NOT NULL AND timezone !=''").Find(&users).Error

	if err != nil {
		return fmt.Errorf("failed to fetch users: %v", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("no users found with valid timezone")
	}

	timezoneGroups := make(map[string][]User)

	for _, user := range users {
		timezoneGroups[user.Timezone] = append(timezoneGroups[user.Timezone], user)
	}

	for tz, users := range timezoneGroups {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			emailErrors = append(emailErrors, fmt.Sprintf("invalid timezone %s: %v", tz, err))
			continue
		}

		for _, user := range users {
			//Check if it's 8 PM in user's timezone
			userLocalTime := time.Now().In(loc)
			// isEightPM := userLocalTime.Hour() == 20 // 20 = 8 PM in 24-hour format

			// if !isEightPM {
			// 	fmt.Printf("⏰ Skipping %s - Current time in %s is %d:00 (waiting for 20:00/8 PM)\n",
			// 		user.Email, tz, userLocalTime.Hour())
			// 	continue
			// }

			fmt.Printf("🕗 It's 8 PM in %s timezone for user %s - sending email now!\n", tz, user.Email)

			startOfDayLocal := time.Date(userLocalTime.Year(), userLocalTime.Month(), userLocalTime.Day(), 0, 0, 0, 0, loc)
			endOfDayLocal := startOfDayLocal.Add(24 * time.Hour)
			startUTC := startOfDayLocal.UTC()
			endUTC := endOfDayLocal.UTC()

			var summaries []string
			err := r.db.Table("notes").Select("summary").Where("user_id = ? AND created_at>= ? AND created_at <= ?", user.ID, startUTC, endUTC).Pluck("summary", &summaries).Error

			if err != nil {
				emailErrors = append(emailErrors, fmt.Sprintf("failed to get summaries for user %s: %v", user.Email, err))
				continue
			}

			if len(summaries) > 0 {
				combinedSummary := strings.Join(summaries, "\n\n")
				prompt := fmt.Sprintf(`
You are a personal learning assistant. The user wrote several notes today, but they might forget them if not reinforced. Your job is to help them RECALL and PRACTICE what they wrote by creating an engaging daily review.

Create a summary that:

🧠 **RECALL CHALLENGE** (Test their memory):
Start with 2-3 questions to test if they remember key points from their notes:
- "Do you remember what you wrote about...?"
- "Can you recall the main insight you had about...?"
- "What was the key point you discovered regarding...?"

📝 **YOUR NOTES SUMMARY** (Reinforce their learning):
- Summarize their notes in an engaging, easy-to-remember way
- Highlight the most important insights they wrote
- Connect different notes to show patterns in their thinking
- Use their own words and concepts when possible

💡 **KEY TAKEAWAYS TO REMEMBER**:
- Extract 2-3 main lessons from their notes
- Present them as memorable, practical insights
- Help them see the value in what they wrote
🎯 **PRACTICE REMINDER**:
End with encouragement to keep writing and a simple reminder of why their notes matter.

Make it personal, engaging, and focused on helping them remember and value what they wrote today.

User's note summaries from today:
%s
Create their recall-practice summary:`, combinedSummary)
				res, err := r.llm.GenerateNoteSummary(prompt)
				if err != nil {
					finalMessage = "We couldn't generate your summary today, but keep up the great work! 🌟"
				} else {
					finalMessage = res
				}
			} else {
				finalMessage = `📝 Your Daily Reflection

Today was quiet on the note-taking front, but that's perfectly okay!

🌟 Tomorrow is a fresh opportunity to capture your thoughts, ideas, and discoveries.

Keep growing, keep learning! ✨`
			}

			// Send email and track errors
			err = r.sendDailySummaryEmail(user.Email, user.Name, finalMessage)
			if err != nil {
				emailErrors = append(emailErrors, fmt.Sprintf("failed to send email to %s: %v", user.Email, err))
				fmt.Printf("❌ Failed to send email to %s: %v\n", user.Email, err)
			} else {
				sentCount++
				fmt.Printf("✅ Email sent successfully to %s at 8 PM %s time\n", user.Email, tz)
			}
		}
	}

	// Print summary
	totalUsers := len(users)
	skippedUsers := totalUsers - sentCount - len(emailErrors)
	fmt.Printf("📊 Summary: Total users: %d, Sent: %d, Failed: %d, Skipped (wrong time): %d\n",
		totalUsers, sentCount, len(emailErrors), skippedUsers)

	// Return error summary if any emails failed
	if len(emailErrors) > 0 {
		return fmt.Errorf("some emails failed to send:\n%s", strings.Join(emailErrors, "\n"))
	}

	return nil
}

func (r *authRepo) sendDailySummaryEmail(recipientEmail, recipientName, message string) error {
	mjApiKeyPublic := os.Getenv("MJ_APIKEY_PUBLIC")
	mjApiKeyPrivate := os.Getenv("MJ_APIKEY_PRIVATE")

	if mjApiKeyPublic == "" || mjApiKeyPrivate == "" {
		return fmt.Errorf("MJ_APIKEY_PUBLIC and MJ_APIKEY_PRIVATE must be set in environment")
	}

	fmt.Printf("🔧 Sending email via Mailjet to: %s\n", recipientEmail)
	fmt.Printf("🔑 Using API Key: %s...\n", mjApiKeyPublic[:10])

	// Create Mailjet client
	mailjetClient := mailjet.NewMailjetClient(mjApiKeyPublic, mjApiKeyPrivate)

	// Create the email message
	subject := "Your Daily Note Summary 📝"
	htmlContent := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Your Daily Summary</title>
			<style>
				body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
				.content { background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; }
				.message { background: white; padding: 20px; border-radius: 8px; margin: 20px 0; border-left: 4px solid #667eea; }
				.footer { text-align: center; margin-top: 30px; color: #666; font-size: 14px; }
			</style>
		</head>
		<body>
			<div class="header">
				<h1>📝 Your Daily NoteMind Summary</h1>
				<p>Hello %s! Here's your personalized daily reflection.</p>
			</div>
			<div class="content">
				<div class="message">
					%s
				</div>
				<p><strong>Keep up the great work!</strong> 🌟</p>
			</div>
			<div class="footer">
				<p>This email was sent by NoteMind - Your Personal Learning Assistant</p>
			</div>
		</body>
		</html>
	`, recipientName, message)

	// Prepare Mailjet message structure
	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: "notereviveapp@gmail.com", // Your verified Mailjet sender email
				Name:  "NoteRevive",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: recipientEmail,
					Name:  recipientName,
				},
			},
			Subject:  subject,
			TextPart: message, // Plain text version
			HTMLPart: htmlContent,
		},
	}

	// Send the email
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	} 

	fmt.Printf("📬 Mailjet response: %+v\n", res)
	fmt.Printf("✅ Email sent successfully to %s\n", recipientEmail)

	return nil
}