package util

type RuleConf struct {
	RuleId      int
	RuleName    string
	RuleContent string
	RuleVersion string

	RuleRunFuncsMap map[string]string
	RuleRunObjsMap  map[string]string
}
