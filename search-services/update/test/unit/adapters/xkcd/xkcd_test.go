package xkcd_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"yadro.com/course/update/adapters/xkcd"
	"yadro.com/course/update/core"
)

func Test_GetLast(t *testing.T) {
	log := slog.Default()
	client, err := xkcd.NewClient("https://xkcd.com", time.Second*10, log)
	assert.NoError(t, err)

	ans, err := client.LastID(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, ans, 3065) // TODO ручками каждый раз тест править плохо, нужно подумоть
}

func Test_GetSpecified(t *testing.T) {
	log := slog.Default()
	client, err := xkcd.NewClient("https://xkcd.com", time.Second*10, log)
	assert.NoError(t, err)

	ans, err := client.Get(context.Background(), 614)
	assert.NoError(t, err)

	assert.Equal(t, ans, core.XKCDInfo{
		ID:          614,
		Title:       "Woodpecker",
		URL:         "https://xkcd.com/614/info.0.json",
		Description: "[[A man with a beret and a woman are standing on a boardwalk, leaning on a handrail.]]\nMan: A woodpecker!\n\u003C\u003CPop pop pop\u003E\u003E\nWoman: Yup.\n\n[[The woodpecker is banging its head against a tree.]]\nWoman: He hatched about this time last year.\n\u003C\u003CPop pop pop pop\u003E\u003E\n\n[[The woman walks away.  The man is still standing at the handrail.]]\n\nMan: ... woodpecker?\nMan: It's your birthday!\n\nMan: Did you know?\n\nMan: Did... did nobody tell you?\n\n[[The man stands, looking.]]\n\n[[The man walks away.]]\n\n[[There is a tree.]]\n\n[[The man approaches the tree with a present in a box, tied up with ribbon.]]\n\n[[The man sets the present down at the base of the tree and looks up.]]\n\n[[The man walks away.]]\n\n[[The present is sitting at the bottom of the tree.]]\n\n[[The woodpecker looks down at the present.]]\n\n[[The woodpecker sits on the present.]]\n\n[[The woodpecker pulls on the ribbon tying the present closed.]]\n\n((full width panel))\n[[The woodpecker is flying, with an electric drill dangling from its feet, held by the cord.]]\n\n{{Title text: If you don't have an extension cord I can get that too.  Because we're friends!  Right?}}",
	})
}
