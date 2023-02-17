package gitcomm

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	bb "github.com/jiananwen1/gitcomm/promptui"
)

var (
	msg = Message{
		Type:     "feat",
		Subject:  "",
		TapdType: "story",
		TapdId:   0,
	}
	//types = []string{
	//	"feat	[new feature]",
	//	"fix		[bug fix]",
	//	"docs	[changes to documentation]",
	//	"style	[format, missing semi colons, etc; no code change]",
	//	"refactor	[refactor production code]",
	//	"test	[add missing tests, refactor tests; no production code change]",
	//	"chore	[update grunt tasks etc; no production code change]",
	//	"other	[user define any content message]",
	//}
	tapdTypes = []string{"story", "bug"}

	chooseLastTapdId = []string{"true", "false"}

	types = []string{
		"feat	[新功能（feature）]",
		"fix		[修补bug]",
		"docs	[文档（documentation）]",
		"style	[格式（不影响代码运行的变动）]",
		"refactor	[重构（即不是新增功能，也不是修改bug的代码变动）]",
		"test	[增加测试]",
		"chore	[构建过程或辅助工具的变动]",
		"other	[自定义消息]",
	}
)

func fillMessage(msg *Message) {
	var err error
	msg.Type, err = bb.PromptAfterSelect("Choose a commit type", types)
	checkInterrupt(err)

	p := bb.Prompt{
		BasicPrompt: bb.BasicPrompt{
			Label:     msg.Type + ":",
			Formatter: linterSubject,
			Validate:  validateSubject,
		},
	}
	msg.Subject, err = p.Run()
	checkInterrupt(err)

	if msg.Type == "other" {
		return
	}
	//	mlBody := bb.MultilinePrompt{
	//		BasicPrompt: bb.BasicPrompt{
	//			Label: "Type in the body",
	//			Default: `# If applied, this commit will...
	//# [Add/Fix/Remove/Update/Refactor/Document] [issue #id] [summary]
	//`,
	//			Formatter: linterBody,
	//			Validate:  validateBody,
	//		},
	//	}
	//	msg.Body, err = mlBody.Run()
	//	checkInterrupt(err)
	// # Why is it necessary? (Bug fix, feature, improvements?)
	// -
	// # How does the change address the issue?
	// -
	// # What side effects does this change have?
	// -
	msg.TapdType, err = bb.PromptAfterSelect("Choose a TAPD type", tapdTypes)
	checkInterrupt(err)

	usingLastStoryId, err := bb.PromptAfterSelect("Using last story id?", chooseLastTapdId)
	checkInterrupt(err)

	b, err := strconv.ParseBool(usingLastStoryId)
	checkInterrupt(err)
	if b {
		lastCommit, err := gitWithOutput("log", "-1")
		checkInterrupt(err)
		tapdId := getStoryIdFromLastCommitLog(lastCommit)
		if len(tapdId) == 0 {
			checkInterrupt(errors.New("invalid story id"))
		}
		msg.TapdId, err = strconv.Atoi(tapdId)
		checkInterrupt(err)
		return
	}

	p = bb.Prompt{
		BasicPrompt: bb.BasicPrompt{
			Label: "--" + msg.TapdType + "=",

			Formatter: linterTapdId,
			Validate:  validateTpadId,
		},
	}
	tapdId, err := p.Run()
	checkInterrupt(err)
	msg.TapdId, err = strconv.Atoi(tapdId)
	checkInterrupt(err)
}

// Prompt function assignes user input to Message struct
func Prompt() string {
	fillMessage(&msg)
	gitMsg := msg.String() + "\n"

	// todo: save git message, tapd id into somewhere
	//Info("\nCommit message is:\n%s", gitMsg)
	//for {
	//	cp := bb.ConfirmPrompt{
	//		BasicPrompt: bb.BasicPrompt{
	//			Label:   "Is everything OK? Continue",
	//			Default: "N",
	//			NoIcons: true,
	//		},
	//		ConfirmOpt: "e",
	//	}
	//	c, err := cp.Run()
	//	checkConfirmStatus(c, err)
	//	if c == "Y" {
	//		break
	//	}
	//	if c == "E" {
	//		numlines := len(strings.Split(gitMsg, "\n"))
	//		for ; numlines > -1; numlines-- {
	//			fmt.Print(bb.ClearUpLine())
	//		}
	//		gitMsg, err = bb.Editor("", gitMsg)
	//		checkInterrupt(err)
	//		Info(gitMsg)
	//		// checkConfirmStatus(bb.Confirm("Is everything OK? Continue", "N", true))
	//		// return gitMsg
	//		continue
	//	}
	//}
	return gitMsg
}

// TagPrompt prompting tag version level to upgrade
func TagPrompt() string {
	s := bb.Select{
		Label: "Choose tag level",
		Items: []string{
			"patch",
			"minor",
			"major",
		},
		Default: 0,
	}
	_, level, err := s.Run()
	checkInterrupt(err)
	return level
}

// PromptConfirm is a common function to ask confirm before some action
func PromptConfirm(msg string) bool {
	c, err := bb.Confirm(msg, "N", false)
	checkInterrupt(err)
	if c == "N" {
		return false
	}
	return true
}

func linterSubject(s string) string {
	if len(s) == 0 {
		return s
	}
	// Remove all leading and trailing white spaces
	s = strings.TrimSpace(s)
	// Remove trailing period ...
	s = strings.TrimSuffix(s, "...")
	// Then strings.Title the first word in string
	flds := strings.Fields(s)
	flds[0] = strings.Title(flds[0])
	return strings.Join(flds, " ")
}
func linterTapdId(s string) string {
	if len(s) == 0 {
		return s
	}
	// Remove all leading and trailing white spaces
	s = strings.TrimSpace(s)
	return s
}

func linterBody(s string) string {
	if len(s) == 0 {
		return s
	}
	// remove all leading white space
	// doesn't work because there is commented message goes first
	// s = strings.TrimLeft(s, "\t\n\v\f\r")
	var upl = func(sl string) string {
		rs := []rune(sl)
		if len(rs) > 0 {
			rs[0] = unicode.ToUpper(rs[0])
		}
		return string(rs)
	}
	out := []string{}
	lines := strings.Split(s, "\n")
	for i := range lines {
		// if the line is commented with # at the start pass that line
		if len(lines[i]) > 0 && lines[i][0] == '#' {
			continue
		}
		if len(lines[i]) > 72 {
			nl := wrapLine(lines[i], 72)
			out = append(out, nl...)
			continue
		}
		out = append(out, lines[i])
	}
	for {
		if len(out) == 0 {
			break
		}
		if strings.TrimSpace(out[0]) != "" {
			out[0] = upl(out[0])
			break
		}
		out = out[1:]
	}
	return strings.TrimSpace(strings.Join(out, "\n"))
}

func wrapLine(l string, n int) []string {
	words := strings.Split(l, " ")
	out := []string{}
	s := ""
	i := 0
	for {
		if i >= len(words) {
			out = append(out, s)
			break
		}
		separator := " "
		if s == "" {
			separator = ""
		}
		sv := s + separator + words[i]
		if len(sv) >= n {
			out = append(out, s)
			sv = ""
			i--
		}
		s = sv

		i++
	}
	return out
}

func linterFoot(s string) string {
	if len(s) == 0 {
		return s
	}
	s = strings.TrimSpace(s)
	// Split string to lines
	strs := strings.Split(s, "\n")
	out := []string{}
	for i := len(strs); i > 0; i-- {
		if strings.TrimSpace(strs[i-1]) == "" {
			continue
		}
		if strings.HasPrefix(strs[i-1], "* ") {
			strs[i-1] = strings.TrimPrefix(strs[i-1], "* ")
		}
		strs[i-1] = linterSubject(strs[i-1])
		strs[i-1] = "* " + strs[i-1]
		out = append(append([]string{}, strs[i-1]), out...)
	}
	return strings.Join(out, "\n")
}

func validator(n int) func(val interface{}) error {
	return func(val interface{}) error {
		// since we are validating an Input, the assertion will always succeed
		if str, ok := val.(string); !ok || str == "" || len(str) > n {
			return fmt.Errorf("This response cannot be longer than %d characters", n)
		}
		return nil
	}
}

func validateBody(s string) error {
	if s == "" {
		return bb.NewValidationError("Body must not be empty string")
	}
	ins := strings.Split(s, "\n")
	for i := range ins {
		if len(ins[i]) > 72 {
			return bb.NewValidationError("Body must be wraped at 72 characters")
		}
	}
	return nil
}

func validateSubject(s string) error {
	if s == "" {
		return bb.NewValidationError("Subject must not be empty string")
	}
	if len(s) > 72 {
		return bb.NewValidationError("Subject cannot be longer than 72 characters")
	}
	return nil
}

func validateTpadId(s string) error {
	if s == "" {
		return bb.NewValidationError("tpad id must not be empty")
	}
	_, err := strconv.Atoi(s)
	if err != nil {
		return bb.NewValidationError("Tapd id must be a number")
	}
	return nil
}

func checkInterrupt(err error) {
	if err != nil {
		if err != bb.ErrInterrupt {
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
		}
		log.Printf("interrupted by user\n")
		os.Exit(1)
	}
}

func checkConfirmStatus(c string, err error) {
	checkInterrupt(err)
	if c == "N" {
		log.Printf("Commit flow interrupted by user\n")
		os.Exit(1)
	}
}
