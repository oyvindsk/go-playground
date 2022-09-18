
Test SQLite on Google Cloud Run.

The purpose is to run low traffic web apps for (almost?) free on Cloud Run, while beeing easily portable (no properitary GCP databases).

Basically: 
 - DB stored "locally" on the Cloud Run instance, on disk, but AFAIK disk == memory
 - Stream DB to Storage with Litestream
 - Read DB on instance start
 - Limit Cloud Run to 0 or 1 instance 

 - Too much data: Copying to/from Storage will become slow and expensinve
 - Too much traffic: If it never scales to 0 a VM will be cheaper and simpler
 - Too much traffic: There are limits to what 1 process + SQLite can do.


SQLite: 

Go and packages:

Google Cloud Platform:
project=sqlite-test-353918
$ gcloud config set account oyvindska@gmail.com
$ gcloud auth application-default login


Litestream:
 - Compiling:       https://litestream.io/install/source/
    - git checkout v0.3.8
    - go install ./cmd/litestream
 - SQLite tips:     https://litestream.io/tips/
 - Sync to GCS:     https://litestream.io/guides/gcs/

Replicate:
$ litestream replicate ./foo.db 'gcs://oyvindsk-sqlite-test-litestream'

Restore:
$ DIR=foo2-backup-$(date '+%s%N')
$ mkdir $DIR
$ mv foo2.db* $DIR
$ litestream restore -o ./foo2.db  'gcs://oyvindsk-sqlite-test-litestream'



https://mermaid.live/edit#pako:eNqNkkFrhDAQhf_KkFMXdlnUm4eWslt6aaHU3tRDMGMVYmJj0rKI_72j0e26vexlCLzvvclM0rNCC2Qx-zS8reDjCACZogK1CNL98-ENEsQOrOFlWRf7HHa7eygCzxTBXXrQyvJaoYHEcmPzzUTIIH2pLXbWIG_gnQ7aYO5dMpiQhlzpKxV4bFvqYr6RwC9HbLcmZXQZlpxUMetjxEQInYigTypnhf5RD8Nsj85i-E-cLKSD0ufr_CnhhSKjK8eJFrJYEuFnWOLzq4wbWC_5VuFq0jUnQ7_8cL103eYbtmUNGgoS9JT9iGfMVthgxmI6Ciy5kzZjmRoIda3gFp9ETY_C4pLLDreMO6vH1bLYGocLdKw5_YxmpoZfnNmwcw

graph TD   
    id1[/GCP Sees traffic/] --> c1
    c1([Container Start]) --> l1[Litestream Restore]
    l1 --> main[Main App Serve Requests]
    l1 --> l3[Litestream Sync]
    main --> doSd1{Shutdown?}
    l3 --> doSd2{Shutdown?}
    doSd1 -- no --> main
    doSd2 -- no --> l3
    doSd1 -- yes --> mainSd[Main Shutdown]
    doSd2 -- yes --> mainSd[Main Shutdown]
    mainSd --> l2[Litestream Shutdown]
    l2 --> c2([Container Stop])



###################3
Running: 

Create the config file SECRET-config.sh , see Config-example.sh

Create the Cloud Artifact Registry Repository named in the config file
TODO: Put this in some script? 


Building with GCP build => Artifact Registry,
   and 
Deploying from GCP Artifact Registry => Cloud Run
$ ./scripts/gcp/build-to-artifact-registry-and-deploy-to-cloud-run.sh



Building with Docker locally:
$ ./scripts/locally/build-with-docker-and-run.sh

Configs and secrets:


Plan:
[x] Simple templates (no css)
[x] Basic functionality with htmx
[ ] Support password
[ ] Support message deletion
[ ] Handle errors :P
[ ] Reorganize repo ond go code
[ ] Good sqlite usage
[x] Better templates - css :)
