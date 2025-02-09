package quotes

import (
	"math/rand"
)

// Collection of wisdom quotes
var quotes = []string{
	"The only true wisdom is in knowing you know nothing. - Socrates",
	"The journey of a thousand miles begins with one step. - Lao Tzu",
	"Know thyself. - Ancient Greek Aphorism",
	"The unexamined life is not worth living. - Socrates",
	"To be yourself in a world that is trying to make you something else is the greatest accomplishment. - Ralph Waldo Emerson",
	"The only way to do great work is to love what you do. - Steve Jobs",
	"Wisdom comes from experience, and experience comes from mistakes. - Unknown",
	"The best revenge is to be unlike him who performed the injury. - Marcus Aurelius",
	"Everything has beauty, but not everyone sees it. - Confucius",
	"The more that you read, the more things you will know. - Dr. Seuss",
}

// GetRandomQuote returns a random quote from the collection
func GetRandomQuote() string {
	return quotes[rand.Intn(len(quotes))]
}
