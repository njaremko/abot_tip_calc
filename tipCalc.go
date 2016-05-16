package tip

import (
	"log"
	"strconv"
	"strings"

	"github.com/itsabot/abot/shared/datatypes"
	"github.com/itsabot/abot/shared/language"
	"github.com/itsabot/abot/shared/nlp"
	"github.com/itsabot/abot/shared/plugin"
)

var p *dt.Plugin

func init() {

	trigger := &nlp.StructuredInput{
		Commands: []string{"what", "how", "calculate"},
		Objects:  []string{"tip"},
	}

	fns := &dt.PluginFns{Run: Run, FollowUp: FollowUp}

	var err error
	pluginPath := "github.com/njaremko/abot_tip_calc"
	p, err = plugin.New(pluginPath, trigger, fns)
	if err != nil {
		log.Fatalln("building", err)
	}

	p.Vocab = dt.NewVocab(
		dt.VocabHandler{
			Fn:      parseTip,
			Trigger: trigger,
		},
	)
}

func Run(in *dt.Msg) (string, error) {
	return FollowUp(in)
}

func FollowUp(in *dt.Msg) (string, error) {
	return p.Vocab.HandleKeywords(in), nil
}

func parseTip(in *dt.Msg) string {
	// Amount of money spent
	var amount float64 = 0
	// Default to a 15% tip to be safe
	var tip float64 = 15
	// Input sentence separated by
	var tokenizedSentence []string = in.Tokens

	for i := 0; i < len(tokenizedSentence); i++ {
		// Handle specified percentage for tips
		if strings.Contains(tokenizedSentence[i], "%") {
			if tempTip, err := strconv.ParseFloat(strings.TrimSuffix(tokenizedSentence[i], "%"), 64); err == nil {
				tip = tempTip
			}
		} else if i < len(tokenizedSentence)-1 && tokenizedSentence[i+1] == "percent" {
			if tempTip, err := strconv.ParseFloat(tokenizedSentence[i], 64); err == nil {
				tip = tempTip
			}
		} else {
			// If the previous two cases aren't true, but this is. Then it should be the amount spent.
			val, err := language.ExtractCurrency(tokenizedSentence[i])
			if err == nil {
				amount = float64(val)
			}
		}
	}
	// Return the final string
	if amount != 0 {
		return "I recommend you tip $" + calcTip(amount, tip) + "."
	} else {
		return "I'm sorry, but you didn't specify an amount of money."
	}
}

func calcTip(spent float64, tip float64) string {
	result := spent * tip / 10000
	return strconv.FormatFloat(result, 'f', 2, 64)
}
