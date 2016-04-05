package tip

import (
	"github.com/itsabot/abot/shared/datatypes"
	"github.com/itsabot/abot/shared/nlp"
	"github.com/itsabot/abot/shared/plugin"
	"log"
	"strconv"
	"strings"
)

var p *dt.Plugin

func init() {

	trigger := &nlp.StructuredInput{
		Commands: []string{"what"},
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
			Fn: parseTip,
			Trigger: &nlp.StructuredInput{
				Commands: []string{"what"},
				Objects:  []string{"tip"},
			},
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
	var amount float64 = 0
	var tip float64 = 15

	for _, obj := range in.Tokens {
		if strings.Contains(obj, "$") {
			amount, _ = strconv.ParseFloat(strings.TrimPrefix(obj, "$"), 64)
		}

		if strings.Contains(obj, "%") {
			tip, _ = strconv.ParseFloat(strings.TrimSuffix(obj, "%"), 64)
		}
	}
	if amount != 0 {
		return "I recommend you tip $" + calcTip(amount, tip)
	} else {
		return "Please make sure you include the '$' symbol prefix."
	}
}

func calcTip(spent float64, tip float64) string {
	result := spent * tip / 100
	return strconv.FormatFloat(result, 'f', 2, 64)
}
