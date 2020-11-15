package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

func main() {
	// cluster := gocql.NewCluster("localhost:9042")
	// cluster.Keyspace = "test"
	// cluster.Consistency = gocql.LocalOne
	// session, _ := cluster.CreateSession()
	// defer session.Close()

	// res, err := doCheckRuleIDUniqueness(context.TODO(), session, "1")
	// fmt.Println("1", res)
	// fmt.Println("error", err)

	// res2, err := doCheckRuleIDUniqueness(context.TODO(), session, "voucher_adidas_1")
	// fmt.Println("voucher_adidas_1", res2)
	// fmt.Println("error", err)

	// ruleConfigs := []RuleConfig{
	// 	{
	// 		Rank: 1,
	// 		Condition: []RuleConfigCondition{
	// 			{
	// 				Key:   "spending_habbit",
	// 				Value: "top_spender",
	// 			},
	// 		},
	// 		Result: RuleConfigResult{
	// 			ResultType: "answer",
	// 			ResultDetail: []RuleConfigResultResultDetail{
	// 				{
	// 					Key:   "promo_section_id",
	// 					Value: "adidas_promo",
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	// pr := NewPersonalizationRule(
	// 	"voucher_adidas_2",
	// 	"pnr",
	// 	"Rule for vourcher adidas 1",
	// 	"active",
	// 	"batch",
	// 	"user_label",
	// 	ruleConfigs,
	// 	"",
	// 	uint64(123),
	// 	"putra",
	// 	time.Now(),
	// )

	// err := doCreatePersonalizationRule(context.TODO(), session, pr)
	// if err != nil {
	// 	log.Fatal("errorr ", err)
	// }

	// fmt.Println("success")

	Cassandra2()

}

const checkRuleIDUniquenessStmt = `SELECT rule_id 
																		FROM personalization_rules 
																		WHERE rule_id = ?`

func doCheckRuleIDUniqueness(ctx context.Context, cassandra *gocql.Session, ruleID string) (bool, error) {
	var ruleIDRes string
	err := cassandra.Query(checkRuleIDUniquenessStmt, ruleID).WithContext(ctx).Scan(&ruleIDRes)

	if err == gocql.ErrNotFound {
		return true, nil
	}

	if err != nil {
		return false, err
	}

	return false, nil
}

func NewPersonalizationRule(
	ruleID string,
	ruleOwner string,
	ruleDescription string,
	ruleStatus string,
	ruleProcessingType string,
	ruleType string,
	ruleConfig []RuleConfig,
	ruleBigqueryRawSQL string,
	ruleLastUpdateUserID uint64,
	ruleLastUpdateUserName string,
	ruleLastUpdateAt time.Time,
) PersonalizationRule {
	personalizationRule := PersonalizationRule{
		RuleID:                 ruleID,
		RuleOwner:              ruleOwner,
		RuleDescription:        ruleDescription,
		RuleStatus:             ruleStatus,
		RuleProcessingType:     ruleProcessingType,
		RuleType:               ruleType,
		RuleConfig:             ruleConfig,
		RuleBigqueryRawSQL:     ruleBigqueryRawSQL,
		RuleLastUpdateUserID:   ruleLastUpdateUserID,
		RuleLastUpdateUserName: ruleLastUpdateUserName,
		RuleLastUpdateAt:       ruleLastUpdateAt,
	}

	return personalizationRule
}

type PersonalizationRule struct {
	RuleID                 string
	RuleOwner              string
	RuleDescription        string
	RuleStatus             string
	RuleProcessingType     string
	RuleType               string
	RuleConfig             []RuleConfig
	RuleBigqueryRawSQL     string
	RuleLastUpdateUserID   uint64
	RuleLastUpdateUserName string
	RuleLastUpdateAt       time.Time
}

// RuleConfig is struct for Personalization Rule Config
type RuleConfig struct {
	Rank      int                   `json:"rank"`
	Condition []RuleConfigCondition `json:"condition"`
	Result    RuleConfigResult      `json:"result"`
}

// RuleConfigCondition is struct for Personalization Rule Config Condition
type RuleConfigCondition struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// RuleConfigResult is struct for Personalization Rule Config Result
type RuleConfigResult struct {
	ResultType   string                         `json:"result_type"`
	ResultDetail []RuleConfigResultResultDetail `json:"result_detail"`
}

// RuleConfigResultResultDetail is struc for Personalization Rule Config Result Detail
type RuleConfigResultResultDetail struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

const createPersonalizationRulesStmt = `INSERT INTO personalization_rules 
																					(rule_id, rule_owner, rule_description, rule_status, rule_processing_type, rule_type, rule_config, rule_bigquery_raw_sql,
																						rule_last_update_user_id, rule_last_update_user_name, rule_last_update_at)
																				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

func doCreatePersonalizationRule(context context.Context, cassandra *gocql.Session, personalizationRule PersonalizationRule) error {
	JSONRuleConfigs, err := json.Marshal(personalizationRule.RuleConfig)
	if err != nil {
		return err
	}

	err = cassandra.Query(createPersonalizationRulesStmt,
		personalizationRule.RuleID,
		personalizationRule.RuleOwner,
		personalizationRule.RuleDescription,
		personalizationRule.RuleStatus,
		personalizationRule.RuleProcessingType,
		personalizationRule.RuleType,
		string(JSONRuleConfigs),
		personalizationRule.RuleBigqueryRawSQL,
		personalizationRule.RuleLastUpdateUserID,
		personalizationRule.RuleLastUpdateUserName,
		personalizationRule.RuleLastUpdateAt,
	).WithContext(context).Exec()

	return err
}

// type Name struct {
// 	FirstName string `json:"first_name"`
// 	LastName  string `json:"last_name"`
// 	Age       int    `json:"age"`
// }
type Name struct {
	FirstName string
	LastName  string
	Age       int
}

func Cassandra2() {
	cluster := gocql.NewCluster("localhost:9042")
	cluster.Keyspace = "test"
	cluster.Consistency = gocql.LocalOne
	session, _ := cluster.CreateSession()
	defer session.Close()

	name := []Name{
		{
			FirstName: "putri",
			LastName:  "mutiara",
			Age:       20,
		},
	}

	nameJSON, err := json.Marshal(name)
	if err != nil {
		log.Fatal(err)
	}

	id := "4"
	nameString := string(nameJSON)

	if err := session.Query(`INSERT INTO tests (id, name) VALUES (?,?)`, id, nameString).Exec(); err != nil {
		log.Fatal(err)
	}

	var idRes string
	var nameRes string

	if err := session.Query(`SELECT id, name FROM tests WHERE id = ?`, id).Scan(&idRes, &nameRes); err != nil {
		if err == gocql.ErrNotFound {
			log.Fatal("errorrr select test not founddd.....")
		}

		log.Fatal("error lain ", err)
	}
	fmt.Println("id", idRes)
	fmt.Println("Name", nameRes)

	nameResJSON := []Name{}
	err = json.Unmarshal([]byte(nameRes), &nameResJSON)
	if err != nil {
		log.Fatal("errorrr ", err)
	}

	fmt.Println(nameResJSON)
}

func Cassandra1() {
	// connect to the cluster
	cluster := gocql.NewCluster("localhost:9042")
	cluster.Keyspace = "test"
	cluster.Consistency = gocql.LocalOne
	session, _ := cluster.CreateSession()
	defer session.Close()

	// insert a tweet
	// if err := session.Query(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`,
	// 	"me", gocql.TimeUUID(), "hello world").Exec(); err != nil {
	// 	log.Fatal(err)
	// }

	var id gocql.UUID
	var text string

	/* Search for a specific set of records whose 'timeline' column matches
	 * the value 'me'. The secondary index that we created earlier will be
	 * used for optimizing the search */
	if err := session.Query(`SELECT id, text FROM tweet WHERE text = ? LIMIT 1`,
		"hello world").Consistency(gocql.One).Scan(&id, &text); err != nil {
		log.Fatal("errorrrr ", err)
	}
	fmt.Println("Tweet:", id, text)

	// list all tweets
	iter := session.Query(`SELECT id, text FROM tweet WHERE timeline = ?`, "me").Iter()
	for iter.Scan(&id, &text) {
		fmt.Println("Tweet:", id, text)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
}
