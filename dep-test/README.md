
From gopher slack/#vendor:

I have an issue with `dep init -gopath` pulling in the wrong versions of transitive dependencies.
wondering if it could be a bug or if there's something I'm misunderstanding
Basically it sometimes uses the tip from github instead of what I have i $GOPATH:


Shouldn't the revision of x/net/ be what I already have in $GOPATH?: 

    os@oslapto ~$ cd $GOPATH/src/github.com/oyvindsk/go-playground/dep-test/
    os@oslapto (master) ~/go/src/github.com/oyvindsk/go-playground/dep-test$ dep init -gopath
    Searching GOPATH for projects...
    Using master as constraint for direct dep cloud.google.com/go
    Locking in master (290422c) for direct dep cloud.google.com/go
    Locking in v1.0.0 (150dc57) for transitive dep google.golang.org/appengine
    Locking in master (921ae39) for transitive dep golang.org/x/oauth2
    Locking in master (1d60e46) for transitive dep golang.org/x/sync
    Locking in v1.0.0 (9255415) for transitive dep github.com/golang/protobuf
    Locking in master (b68f304) for transitive dep golang.org/x/net
    Locking in v1.11.1 (1e2570b) for transitive dep google.golang.org/grpc
    Locking in v0.7.0 (076344b) for transitive dep go.opencensus.io
    Locking in master (35de241) for transitive dep google.golang.org/genproto
    Locking in master (3097bf8) for transitive dep google.golang.org/api
    Locking in v0.3.0 (f21a4df) for transitive dep golang.org/x/text
    Locking in v2.0.0 (317e000) for transitive dep github.com/googleapis/gax-go

    os@oslapto (master) ~/go/src/github.com/oyvindsk/go-playground/dep-test$ dep status
    PROJECT                       CONSTRAINT     VERSION        REVISION  LATEST   PKGS USED
    cloud.google.com/go           branch master  branch master  290422c   b475e33  5   
    github.com/golang/protobuf    v1.0.0         v1.0.0         9255415   v1.0.0   8   
    github.com/googleapis/gax-go  v2.0.0         v2.0.0         317e000   v2.0.0   1   
    go.opencensus.io              v0.7.0         v0.7.0         076344b   v0.7.0   12  
    golang.org/x/net              branch master  branch master  b68f304   b68f304  8   
    golang.org/x/oauth2           branch master  branch master  921ae39   921ae39  5   
    golang.org/x/sync             branch master  branch master  1d60e46   1d60e46  1   
    golang.org/x/text             v0.3.0         v0.3.0         f21a4df   v0.3.0   14  
    google.golang.org/api         branch master  branch master  3097bf8   3097bf8  8   
    google.golang.org/appengine   v1.0.0         v1.0.0         150dc57   v1.0.0   12  
    google.golang.org/genproto    branch master  branch master  35de241   35de241  9   
    google.golang.org/grpc        v1.11.1        v1.11.1        1e2570b   v1.11.1  24  

    os@oslapto (master) ~/go/src/github.com/oyvindsk/go-playground/dep-test$ cd $GOPATH/src/golang.org/x/net/ && git log -n 1 --pretty=oneline
    0744d001aa8470aaa53df28d32e5ceeb8af9bd70 (HEAD -> master, origin/master, origin/HEAD) proxy: add mention of RFC 1929 for SOCKS5

## Others:

### Glide
Does not support getting versions/revisions from existing GOPATH

### Govendor
Does copy files from the local GOPATH into vendor. Seems like it picks the correct versions:

    govendor init
    govendor add +e
