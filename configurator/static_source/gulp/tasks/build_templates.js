var gulp = require('gulp'),
    conf = require('../config').build_templates,
    connect = require('gulp-connect'),
    templateCache = require('gulp-angular-templatecache');

for (var key in conf) {

    (function (_task, _source, _filename, _dest, _prefix) {

        var task = 'template:' + _task,
            source = _source,
            filename = _filename,
            prefix = _prefix,
            dest = _dest;

        gulp.task(task, function() {
            return gulp.src(source)
                .pipe(templateCache(filename, {
                    transformUrl: function(url) {
                        return prefix + url.replace(/\.html\.html/, '.html')
                    }
                }))
                .pipe(gulp.dest(dest))
                .pipe(connect.reload());
        });


    })(key, conf[key].source, conf[key].filename, conf[key].dest, conf[key].prefix)
}

gulp.task('template', ['template:public', 'template:private']);
