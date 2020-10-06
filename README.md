# OVH GO DNS RECORD CREATION AND DELETE

This very simple script will create or Delete a Subdomain in a given domain on OVH DNS.

I use this script in a Helm chart, at the creation and deletion to create and delete websites in a Kubernetes CLuster

----

# HOWTO

## OVH configuration 
The OVH sdk will first look for environment variables to indentify to the API, if none is found, it will look for a configuration file named ovh.conf.

More info on the authentification proces on the official github: https://github.com/ovh/go-ovh#configuration

## Script configuration

If you have not modified the recordtype var in the code, the script will look for four environment variables, that need to be set.

* *OVH_DOMAIN*: name of the domain you want to create a record to
* *OVH_SUBDOMAIN*: name iof the subdomain you want to create a record to
* *OVH_IP_ENDPOINT*: target of the record
* *OVH_ACTION*: action to take, the can take **CREATE** or **DELETE**
* *OVH_RECORD_TYPE*: type of record, if none is provided, A record will be selected by default

## HOW DO I USE THIS SCRIPT

I created a docker with this script.

Using helm chart hook, I launch a job at the creation of the chart that create a DNS record.

When the chart is deleted, another Job is launched to delete the DNS record.

The arguments (domain,subdomain and IP), are provided in the value.yaml, and passed in configmapand mounted as environment varialble in the job's pod.

For the OVH ID, is didn't manage to make the script work with the id info set as environment variables,  so I created a secret with the ovh.conf file and mount it in the root home folder, where the sdk will look.

## Dockerfile

I provided the dockerfile you can use it if you wish, and add your arg in the environment variable.
