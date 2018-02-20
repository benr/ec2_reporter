package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/olekukonko/tablewriter"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ec2.New(sess)

	instances, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}

	var instance_data [][]string

	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {

			var instance_name string
			var dept_name string
			var contact_name string
			var tags string
			for _, tag := range instance.Tags {
				//fmt.Printf("  ---> TAG: %20s = %s\n", *tag.Key, *tag.Value)
				/*
					if *tag.Key == "Name" {
						instance_name = *tag.Value
					} else {
						tags += *tag.Key + "=" + *tag.Value + " "
					}
				*/
				switch *tag.Key {
				case "Name":
					instance_name = *tag.Value
				case "X-Dept":
					dept_name = *tag.Value
				case "X-Contact":
					contact_name = *tag.Value
				default:
					tags += *tag.Key + "=" + *tag.Value + " "
				}
			}
			//fmt.Printf("%24s: %12s %24s %s [%s]\n", instance_name,
			//	*instance.InstanceId, *instance.InstanceType, *instance.LaunchTime,
			//	tags)
			idata := []string{instance_name, string(*instance.InstanceId),
				string(*instance.InstanceType), dept_name, contact_name}
			instance_data = append(instance_data, idata)
		}
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "ID", "Type", "Dept", "Contact"})
	table.SetAutoWrapText(false)
	for _, v := range instance_data {
		table.Append(v)
	}
	table.Render()
}
