module sigs.k8s.io/legacyflag

go 1.12

// Below require/replace pin the same versions used by k/k.

require github.com/spf13/pflag v1.0.1

replace github.com/spf13/pflag => github.com/spf13/pflag v1.0.1
