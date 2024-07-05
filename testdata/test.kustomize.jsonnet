local project = std.extVar("project");

{
  "$files": {
    "a.txt": "a\n",
  },
  // Writes files and run kustomize-edit-add-resource within the directory
  // to update kustomization.yaml
  "$kustomize": {
    "path/to/kustomization.yaml/dir": {
      ["projects/%(project)s.yaml" % { project: "myproject" }]: |||
        metadata:
          name: "%(project)s"
||| % { project: "myproject" }
    },
  },
}
