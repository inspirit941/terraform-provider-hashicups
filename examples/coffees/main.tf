terraform {
  required_providers {
    hashicups = {
      source = "hashicorp.com/edu/hashicups"
    }
  }
}

// provider block에서, 필요한 attribute를 정의해준다.
// provider 코드에서 '환경변수' 로도 해당 값을 받을 수 있도록 해두었지만, 쉽게 run the cycle하기 위해 여기에 평문으로도 추가.
provider "hashicups" {
  host     = "http://localhost:19090"
  username = "education"
  password = "test123"
}

data "hashicups_coffees" "edu" {}

// verify data sources를 위한 block
output "edu_coffees" {
  value = data.hashicups_coffees.edu
}