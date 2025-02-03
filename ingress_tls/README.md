# AWS loadbalancer controller with AWS Certificate Manager

Here if you see i am using the aws loadbalancer controller for my project

# prblem

Here my major problem is 

i am having my own domain and certifcate
which is crt key and pem file but in kubernete if i need to add i need to give like secretname and i need to create manually secret file
its good but not good like all the time we do it manually

# solution

for that easy way we need to add our certificate to ** aws certificate manager ** and after that we need to add some of the annotation
like
    service.beta.kubernetes.io/aws-load-balancer-internal: "false"
    service.beta.kubernetes.io/aws-load-balancer-healthcheck-protocol: http
    service.beta.kubernetes.io/aws-load-balancer-additional-resource-tags: Environment=test
    service.beta.kubernetes.io/aws-load-balancer-healthcheck-timeout: "30"
    service.beta.kubernetes.io/aws-load-balancer-ssl-cert: arn:aws:acm:us-west-2:xxxxx:certificate/xxxxxxx # Here you need to create ACM and after that it will give arn you need to copy and past here.
    service.beta.kubernetes.io/aws-load-balancer-ssl-ports: "443"
    service.beta.kubernetes.io/aws-load-balancer-ssl-negotiation-policy: ELBSecurityPolicy-TLS13-1-2-2021-06

here the main thin is arn which is ssl-cert 

this we need to copy our arn and add here to confirm workin or not

