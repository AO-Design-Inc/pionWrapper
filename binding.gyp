{
  "targets": [
    {
      "target_name": "pionWrapper",
      "sources": [ "src/pionWrapper.c" ],
      "libraries": [ 
	"<(module_root_dir)/go-src/pionHandler.so"
	]
    }
  ]
}
