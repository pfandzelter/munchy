.PHONY: deploy plan clean

deploy: munchy.zip main.tf init.done
	terraform apply
	touch $@

plan: munchy.zip main.tf init.done
	terraform plan
	touch $@

init.done:
	terraform init
	touch $@

munchy.zip: munchy
	chmod +x munchy
	zip -j $@ $<

munchy: main.go format.go dynamo.go
	go get .
	GOOS=linux GOARCH=amd64 go build -ldflags="-d -s -w" -o $@

clean:
	terraform destroy
	rm -f init.done deploy.done munchy.zip munchy