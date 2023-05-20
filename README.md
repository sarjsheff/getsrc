### GetSrc - git web viewer

Git ui with http dumb clone support.

[https://git.sheff.online]

Sample config getsrc.yaml :

``` yaml
repos:
  repo1: 
    path: "/tmp/repo1"
    description: "Test repo."
  repo2: 
    path: "/tmp/repo2"
cloneurl: "https://git.sheff.online"
title: "Test title"
seo:
  description: "Dev test."
  title: "title"
  sitename: "MyGit"
  custom: |
    <meta property="og:locale" content="ru_RU" />
```

minimal:

``` yaml
repos:
  myrepo: 
    path: "/tmp/mygitrepo"
```

Docker:

```
docker run -it --rm --platform linux/amd64 -v /tmp/gitrepos:/repos -v $(pwd)/getsrc.yaml:/app/getsrc.yaml r.sheff.online/sheff/getsrc:latest 
```
