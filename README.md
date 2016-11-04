[![Build Status](https://travis-ci.org/redforks/css.svg?branch=master)](https://travis-ci.org/redforks/css) [![codebeat badge](https://codebeat.co/badges/d9bcef39-9fac-4c57-8ff8-0b10feb19691)](https://codebeat.co/projects/github-com-redforks-css) [![Go Report Card](https://goreportcard.com/badge/github.com/redforks/css)](https://goreportcard.com/report/github.com/redforks/css) [![Go doc](https://godoc.org/github.com/redforks/css?status.svg)](https://godoc.org/github.com/redforks/css) [![Go Cover](http://gocover.io/_badge/github.com/redforks/css/sprite)](http://gocover.io/github.com/redforks/css/sprite) / [![Go Cover](http://gocover.io/_badge/github.com/redforks/css/writer)](http://gocover.io/github.com/redforks/css/writer)
# CSS Spriter

CSS Spriter is a tool generate sprite by scan .css file, it is better described by an example:

    .icon_object {
      background: url('grp1.object.png');
    }

    .icon_method {
      background: url('grp1.method.png');
    }

    .icon_form {
      background: url('grp1.form.png');
    }

    .icon_list {
      background: url('grp2.list.png');
    }

    .icon_template {
      background: url('grp2.template.png');
    }

    .icon_enum {
      background: url('grp2.enum.png');
    }

    .icon_service {
      background: url('grp2.service.png');
    }

It is a simple .css file defines several icons:

 1. ![](https://github.com/redforks/css/raw/gh-pages/grp1.object.png) grp1.object.png
 1. ![](https://github.com/redforks/css/raw/gh-pages/grp1.method.png) grp1.method.png
 1. ![](https://github.com/redforks/css/raw/gh-pages/grp1.form.png) grp1.form.png
 1. ![](https://github.com/redforks/css/raw/gh-pages/grp2.list.png) grp2.list.png
 1. ![](https://github.com/redforks/css/raw/gh-pages/grp2.template.png) grp2.template.png
 1. ![](https://github.com/redforks/css/raw/gh-pages/grp2.enum.png) grp2.enum.png 
 1. ![](https://github.com/redforks/css/raw/gh-pages/grp2.service.png) grp2.service.png

There are 7 .png files, can be divided into two groups: `grp1` and `grp2` by file name prefix. Now run:

    spriter -i tree.css -o build/out.css

Now you get three files in `build` directory: 
 
 1. `out.css`: rewritten with transformed background css declarations.
 1. `Q-EoXMh-.png`: group `grp1` sprite image ![](https://github.com/redforks/css/raw/gh-pages/Q-EoXMh-.png).
 1. `GyO8rqsS.png`: group `grp2` sprite image ![](https://github.com/redforks/css/raw/gh-pages/GyO8rqsS.png).
 
Content of `out.css`:

    .icon_object {
      background: url(Q-EoXMh-.png) no-repeat;
    }

    .icon_method {
      background: url(Q-EoXMh-.png) no-repeat -16px 0;
    }

    .icon_form {
      background: url(Q-EoXMh-.png) no-repeat -32px 0;
    }

    .icon_list {
      background: url(GyO8rqsS.png) no-repeat;
    }

    .icon_template {
      background: url(GyO8rqsS.png) no-repeat -16px 0;
    }

    .icon_enum {
      background: url(GyO8rqsS.png) no-repeat -32px 0;
    }

    .icon_service {
      background: url(GyO8rqsS.png) no-repeat -48px 0;
    }

`spriter` parses the input .css file, gather all background image files match
the format: `group.name.png`. Group them by `group` name, such as `grp1` and
`grp2` in upper example, then create sprite for each group.

`spriter` compute hash value for each sprite file, and use its prefix for
filename. When sprite file changes, the filename also changes. Perfect for
enables http cache.

### Install

As it is a `Go` application, the easiest way is:

    go get -u github.com/redforks/css/cmd/spriter

`spriter` works well under Linex and OS/X, although not tested, but it should also work under Windows.

### Conclusion

 * Easy to use, put `spriter` into your build process, follow the name pattern and go.
 * Update sprite offset automatically.
 * Http cache safe
 * Contains only referenced images, saves bandwidth.

***

中文介绍，看我的[博客](http://blog.503web.com/default/Spriter---%E8%87%AA%E5%8A%A8%E7%94%9F%E6%88%90-css-sprite/)
