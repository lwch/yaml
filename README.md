# yaml

yaml support include

## render

    str, _ := yaml.Render("test/main.yaml")
    fmt.Println(str)

render result:

    includes:
      #include /home/lwch/src/yaml/test/include.d/*.yaml
      #+++++ /home/lwch/src/yaml/test/include.d/include.yaml +++++
      next:
        #include /home/lwch/src/yaml/test/next.yaml
        #+++++ /home/lwch/src/yaml/test/next.yaml +++++
        title: next
        #----- /home/lwch/src/yaml/test/next.yaml -----
      #----- /home/lwch/src/yaml/test/include.d/include.yaml -----

the relative path will convert to absolute path by current file path

## decode

    var ret struct {
        Includes struct {
            Next struct {
                Title string `yaml:"title"`
            } `yaml:"next"`
        } `yaml:"includes"`
    }
    _ := yaml.Decode("test/main.yaml", &ret)
    fmt.Println(ret)