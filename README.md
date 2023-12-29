# Distributed Controller-Workers system

The app is designed as a framework that can run distributed jobs. Scale the work horizontally.
The controller supplies workers with small batches. and workers verify each item in the batch.
If the verification is successful(customizable) the worker will return the verified item to the controller.

This initial implementation only can brute force using different types of alphabets and verify through the http request. I’m planning to add vocabulary and cryptography in future. But it’s designed to be a generic runner. 

Deployment could be done to managed and unmanaged kubernetes clusters. 
## Set up development
Install [golang](https://go.dev/doc/install)
Install [git](https://github.com/git-guides/install-git)
Checkout the project
Run `go mod -sync`

## Run project locally in dev mode
Both components can be started as stand alone apps from `cmd/worker/main.go` and `cmd/controler/main.go`

## Next steps will describe how to set up kubernetes cluster in AWS cloud. Skip if you have one.

**Create AWS account** 
Details [here](https://portal.aws.amazon.com/gp/aws/developer/registration/index.html) 
Take a note of your access and secret key.

**Install AWS CLI**
Details [here](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)

**Configure AWS CLI**
Run `aws configure`

**Install Terraform CLI** 
Details [here](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)

**Register in Terraform CLoud and create an organization(ORG1)**
Steps [here](https://app.terraform.io/public/signup/account)

**Install Docker and register on Docker Hub** 
How to is [here](https://docs.docker.com/engine/install/)
note your account name(YOUR_DOCKER_ACCOUNT)

**Install kubectl**
How exactly is described [here](https://kubernetes.io/docs/tasks/tools/)

### Create AWS infrastructure

**Following set of steps with create an EKS cluster in AWS cloud**
In order to customize max/min/desired sized for worker and controller and their machine instances edit the following file `dcw/deployments/infra/aws/variables.tf`
There should be one controller but infinite number of workers at least for the brute force. Because it will supply same results. But if you have a different type of job you could scale controller as well. 
Create file `deployments/infra/aws/vars.varfile` with following content:

`access_key="<your_aws_access_key>"`

`secret_key="<your_aws_secret_key>"`

where <..> has to be replaced with valid values(secret and access keys from your AWS account)
This file is added to .gtignore and shouldn't be committed under any circumstances. It's only to set up your local authorization for the Terraform.

Once you've set up and changed all the variables go to `deployments/infra/aws/` (if you are not there yet)
Run following commands
(your organization name goes here)
`export TF_CLOUD_ORGANIZATION=<ORG1>` 
`terraform login`
`terraform init`
`terraform apply -var-file=vars.varfile`
Wait for the thing to create the infrastructure. It will take a while(30 mins aprox)
Once it's done... Same folder. 
*The following command will reconfigure your k8s cluster so back up the config file ~/.kube/config)*
`aws eks --region $(terraform output -raw region) update-kubeconfig --name $(terraform output -raw cluster_name)`
Make sure it's looking at the AWS cluster by running
`kubectl cluster-info`

## Build Docker images
Make sure your code is up to date and your job is configured properly (verify configs in `dcw/configs`)
Make sure Docker is up
From the project root
`docker login`
`docker build --progress=plain --no-cache -f deployments/docker/worker/Dockerfile .`
`docker tag <hash_of_just_built_worker_image> <YOUR_DOCKER_ACCOUNT>/worker`
`docker push <YOUR_DOCKER_ACCOUNT>/worker`

`docker build --progress=plain --no-cache -f deployments/docker/controller/Dockerfile .`
`docker tag <hash_of_just_built_controller_image> <YOUR_DOCKER_ACCOUNT>/controller`
`docker push <YOUR_DOCKER_ACCOUNT>/controller`

## Deploy the application to the k8s cluster
Update the image in  deployments/k8s/controller.yml and deployments/k8s/worker.yml (replace <your_docker_account>/controller)
Assuming you are in the project root
Run from there:
`kubectl apply -f deployments/k8s/controller.yml`
`kubectl apply -f deployments/k8s/cworkerr.yml #`
Make sure pods are running:
`kubectl get pods`
After workes are up the job will start running. Once the result is found workers will stop and there will be "+++++++ !!!!!!!! found result !!!!!! +++++++++" with the result coming next in the controller logs.
To look at the logs:
`kubectl logs get <controller_pod_id>`

## Run containers locally
`docker run  --network host <worker_image_hash>`
`docker run -p 50000:50000 <controller_image_hash`
## Add a new task

Currently only brute force from alphabet is supported as a supplier. But a new one can be added. 
The `TestJob` demonstrates how the components work together.
There is a `test/lambda_handler.py` that is a target of the TestJob. Url is blank intentionally. I don't want you guys to torture my lambda function. But you can deploy [yours](https://docs.aws.amazon.com/lambda/latest/dg/getting-started.html) using this code.
For this particular http request I have following configuration:

`{`
 ` "connAttempts": 10,`\
 ` "connAttemptsTtsSec": 5,`\
 ` "batchSize": 10,`\
 ` "workersFactor": 1,`\
 ` "workersSemaphoreWeight": 2,`\
 ` "verifierConfig": {`\
   ` "method": "POST",`\
   ` "url": "<url_goes_here>",`\
   ` "headers": {`\
   ` "Content-Type": "application/json"`\
   ` },`\
   ` "body": "{\n  \"code\": \"%s\"\n}",`\
   ` "customConfig": {`\
    `  "successStatus": 200`\
   ` }`\
 ` }`\
`}`\
  `connAttempts` - how many times a worker tries to connect to the controller (better to deploy the controller first)
  `connAttemptsTtsSec` - sec to sleep between unsuccessful connect attempts
  `batchSize` - how many items the worker will request from the controller
  `workersFactor` - parallel goroutines per 1 cpu
  `workersSemaphoreWeight` - kind of work cache. Worker will request more work while it's still busy doing the previous batch. So 2 means that it will only request 1 more batch per each cpu * workersFactor. If 3 then it's cache 2 more butches per each goroutine.
  `verifierConfig` - is the main config for the logic
  `method` - http method that will verify the item
  `url` - the url
  `headers` - headers that will be used
  `body` - body with %s where the item received from the controller go
  `customConfig` - defines custom identification of success in this case it'll compare the response status code to what's defined in config(200)

# Only should be used for legal purposes. 

Contributions are welcomed. There is a bunch of stuff to do. 


