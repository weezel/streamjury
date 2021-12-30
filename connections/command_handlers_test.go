package connections

import (
	"fmt"
	"streamjury/gameplay"
	"strings"
	"testing"
)

func Test_handleSubmittedSong(t *testing.T) {
	testString := "esitä pirskaleen kova kappale mistä tulee hieno olo https://napster.com"
	splt := strings.Split(testString, " ")
	fmt.Printf("%s\n", strings.Join(splt[2:len(splt)-1], " "))
}

func Test_getSongDescription(t *testing.T) {
	type args struct {
		message []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{
				message: []string{"levyraati", "esitä", "jätte", "kiva", "biisi!", "https://jättekiva.se"},
			},
			want: "jätte kiva biisi!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSongDescription(tt.args.message); got != tt.want {
				t.Errorf("getSongDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aaa(t *testing.T) {
	players := []gameplay.Player{
		{Name: "p1", Uid: 1},
		{Name: "p2", Uid: 2},
		{Name: "p3", Uid: 3},
	}
	var i int
	var p gameplay.Player

	for range players {
		if p.Name == "p2" {
			break
		}
		i++
		if i < len(players) {
			p = players[i]
		}
	}
	fmt.Println(p)
}
