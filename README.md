# yaml

yaml support include

## render

    str, _ := yaml.Render("test/main.yaml")
    fmt.Println(str)

render result:

    includes:
      #include /home/lwch/src/yaml/test/include.d/*.yaml
      #file: /home/lwch/src/yaml/test/include.d/include.yaml
      next:
        #include /home/lwch/src/yaml/test/next.yaml
        #file: /home/lwch/src/yaml/test/next.yaml
        title: next

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