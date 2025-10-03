package auth

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"notemind/internal/llm"

	"github.com/mailersend/mailersend-go"
	"gorm.io/gorm"
)

type AuthRepo interface {

    Create(user *User) error



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
























func (r *authRepo) SendDailySummary() error{
	var users []User
	var finalMessage string

	err := r.db.Where("timezone IS NOT NULL AND timezone !=''").Find(&users).Error

	if err != nil {
		return err
	}

	timezoneGroups := make(map[string][]User)

	for _, user := range users {
		timezoneGroups[user.Timezone] = append(timezoneGroups[user.Timezone], user)
	}
	for tz, users := range timezoneGroups {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			continue
		}
		localTime := time.Now().In(loc)
		//if localTime.Hour() == 20 {
			for _, user := range users {
				startOfDayLocal := time.Date(localTime.Year(), localTime.Month(), localTime.Day(), 0, 0, 0, 0, loc)
				endOfDayLocal := startOfDayLocal.Add(24 * time.Hour)
				startUTC := startOfDayLocal.UTC()
				endUTC := endOfDayLocal.UTC()

				var summaries []string
				err := r.db.Table("notes").Select("summary").Where("user_id = ? AND created_at>= ? AND created_at <= ?", user.ID, startUTC, endUTC).Pluck("summary", &summaries).Error

				if err != nil {
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
						// ✅ Fixed: Use fallback message
						finalMessage = "We couldn't generate your summary today, but keep up the great work! 🌟"
					} else {
						// ✅ Fixed: Use the generated result
						finalMessage=res 
					}
				} else {
					// ✅ Fixed: Added else block for users with no summaries
					finalMessage = `📝 Your Daily Reflection

Today was quiet on the note-taking front, but that's perfectly okay! 

🌟 Tomorrow is a fresh opportunity to capture your thoughts, ideas, and discoveries.

Keep growing, keep learning! ✨`
				}
				//i think from this line i need to integrate mailer send right
                err = r.sendDailySummaryEmail(user.Email, user.Name, finalMessage)
				if err != nil {
					 fmt.Printf("Failed to send email to %s: %v\n", user.Email, err)
				}
			}
		//}
	}
	return nil 
}

func (r *authRepo) sendDailySummaryEmail(recipientEmail, recipientName, message string) error {
	 apiKey:= os.Getenv("MAILERSEND_API_KEY")
	 if apiKey == "" {
		return fmt.Errorf("MAILERSEND_API_KEY not found in environment")
	}
	ms := mailersend.NewMailersend(apiKey)
	ctx:= context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

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

	from := mailersend.From{
		Name:  "NoteRevive",
		Email: "notereviveapp@gmail.com", // Replace with your verified sender email
	}

	recipients := []mailersend.Recipient{
		{
			Name:  recipientName,
			Email: recipientEmail,
		},
	}

	messageObj := ms.Email.NewMessage()

	messageObj.SetFrom(from)
	messageObj.SetRecipients(recipients)
	messageObj.SetSubject(subject)
	messageObj.SetHTML(htmlContent)
	messageObj.SetTags([]string{"daily-summary", "notemind"})

	_, err := ms.Email.Send(ctx, messageObj)
	return err
    
}