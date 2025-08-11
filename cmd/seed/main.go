package main

import (
	"context"
	"flag"
	"fmt"
	"golang-whatsapp-clone/config"
	"golang-whatsapp-clone/database"
	db "golang-whatsapp-clone/database/gen"

	"github.com/jackc/pgx/v5/pgtype"
)

var userId = 0

// This function will seed our current database
// because its a seed function we'll ignore all non-critical errors
func main() {
	clean := flag.Bool("clean", false, "clean/remove previous data")
	flag.Parse()

	appConfig := config.SetupAppConfig()

	// initialize db
	dbpool := database.SetupDatabase(appConfig.DatabaseURL)

	// sqlc generated queries
	dbQueries := db.New(dbpool)

	ctx := context.Background()

	if *clean {
		fmt.Println("Deleting previous data because of the flag 'clean'... 🗑️")
		_ = dbQueries.RemoveAllUsers(ctx)
		fmt.Println("Sucess to delete previous data ✔️")
	}

	fmt.Println("Starting seed data... ⏳")
	// main user
	mainUser, _ := dbQueries.UpsertUserByGoogleAuthSafe(ctx, db.UpsertUserByGoogleAuthSafeParams{
		Name:      fromStringToPGText("Main User"),
		GoogleID:  fromStringToPGText("google-1"),
		Email:     "main-user-test@gmail.com",
		AvatarUrl: fromStringToPGText(getUserAvatar()),
	})
	fmt.Printf("Main user created: %s ✔️\n", mainUser.Email)

	for i := range 10 {
		user := db.UpsertUserByGoogleAuthSafeParams{
			Name:      fromStringToPGText(fmt.Sprintf("User Test %d", i)),
			GoogleID:  fromStringToPGText(fmt.Sprintf("google-%d", i)),
			Email:     fmt.Sprintf("user-%d-test@fake.com", i),
			AvatarUrl: fromStringToPGText(getUserAvatar()),
		}
		_, _ = dbQueries.UpsertUserByGoogleAuthSafe(ctx, user)
		fmt.Printf("User created: %s ✔️\n", user.Email)
	}

	fmt.Println("Success to seed data... ✅")

}

func getUserAvatar() string {
	userId = userId + 1
	return fmt.Sprintf("https://randomuser.me/api/portraits/men/%d.jpg", userId)
}

func fromStringToPGText(s string) pgtype.Text {
	return pgtype.Text{
		Valid:  true,
		String: s,
	}
}
