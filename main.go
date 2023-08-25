package main

import (
	"fmt"
	"math/rand"
)

type BrokerObjectType string   // example: msgVpn

type IdentifyingAttribute struct {
	key, value string
}

type BrokerObjectAttributes []IdentifyingAttribute  // Described as a set of identifying attributes

// TODO: broker-terraform-provider-generator to generate appropriate data structures
var subtypes = map[BrokerObjectType][]BrokerObjectType{
	"msgVpn":  { "aclProfile", "clientProfile", "queue", "jndiQueue"},
	"queue":   {"subscription"},
}

// Only used in this demo, for real broker instances the name is obtrained from the broker
func getInstanceName(brokerObjectAttributes BrokerObjectAttributes) string {
	instanceNamePrefix := ""
	for i := 0; i < len(brokerObjectAttributes)-1; i++ {
		instanceNamePrefix += brokerObjectAttributes[i].value + "-"
	}
	return instanceNamePrefix + brokerObjectAttributes[len(brokerObjectAttributes)-1].value
}

// TODO: Use the real broker object attributes to generate the TF config
//       - Query object attribute settings from broker
//       - Remove attribute for which the value is set to default
//       - Replace value for attributes that are identifying.
//					Example:
//						For a subscriptionTopic object defined as
//							/msgVpns/{msgVpnName}/queues/{queueName}/subscriptions/{subscriptionTopic}
//						Replace msgVpnName, queueName in the TF config by expressions derived from the parent              
func generateConfig(brokerObjectType BrokerObjectType, brokerObjectAttributes BrokerObjectAttributes) {
	// Query object attributes from broker using SEMP GET
	instanceName := getInstanceName(brokerObjectAttributes)  // only used in this demo
	fmt.Printf("  ## Generated config for %s instance:\n  resource \"%s\" \"%s\"  {}\n\n", instanceName, brokerObjectType, instanceName)
}

// Return the list of instances
//
// TODO:
//       - Query all instances of a BrokerObjectType from the broker
//         - Consider using filters: e.g: "List of all MsgVpn names" at https://docs.solace.com/API-Developer-Online-Ref-Documentation/swagger-ui/software-broker/config/index.html
func getInstances(brokerObjectType BrokerObjectType, parentIdentifyingAttributes BrokerObjectAttributes) []BrokerObjectAttributes {
	var instances []BrokerObjectAttributes
	// Substitute for real query from the broker
	//		Ex: query for a subscriptionTopic objects, where parentIdentifyingAttributes are msgVpnName and queueName
	//		/msgVpns/{msgVpnName}/queues/{queueName}/subscriptions
	nrInstances := rand.Intn(3)+1
	fmt.Printf("# Found %d instances of %s on broker, processing each\n", nrInstances, brokerObjectType)
  for i := 0; i < nrInstances; i++ {
		newInstanceIdentifyingAttributes := parentIdentifyingAttributes
		newInstanceIdentifyingAttribute := IdentifyingAttribute{ key: string(brokerObjectType), value: fmt.Sprintf("%s%d", brokerObjectType, i), }
		newInstanceIdentifyingAttributes = append(newInstanceIdentifyingAttributes, newInstanceIdentifyingAttribute)
		instances = append(instances, newInstanceIdentifyingAttributes)
	}
	return instances
}

// Iterates all instances of a child object
func generateConfigForObjectInstances(brokerObjectType BrokerObjectType, parentIdentifyingAttributes BrokerObjectAttributes) error {
	instances := getInstances(brokerObjectType, parentIdentifyingAttributes)
	for i := 0; i < len(instances); i++ {
		generateConfig(brokerObjectType, instances[i])
		for _, subType := range subtypes[brokerObjectType] {
			fmt.Printf("  Now processing subtype %s\n\n", subType)
			// Will need to pass additional params like the parent name etc. so to construct the appropriate names
			err := generateConfigForObjectInstances(subType, instances[i])
			if err != nil {
				return fmt.Errorf("aborting, run into issues")
			}
		}		
	}
	return nil
}

// TODO:
//        - Embed into exisiting provider code
//          - Add comand-line options, ex: terraform-provider-solacebroker generate -url=https://f93.soltestlab.ca:1943 solace_msg_vpn_rest_delivery_point.my_rdp default/my-rdp my-rdp.tf
//        - Reuse existing client from provider
//          - Ensure initial configuration (url, user, passw, etc.)
func main() {
	
	err := generateConfigForObjectInstances("msgVpn", nil)
	// Another example, do it for all queues under a certain msgVpn
	// err := generateConfigForObjectInstances("queue", BrokerObjectAttributes{
	// 	IdentifyingAttribute{ key: "msgVpn", value: "msgVpn2",},})

	if err != nil {
		fmt.Println(err.Error())
	}
  fmt.Println("\nTerraform generate config completed successfully")
}