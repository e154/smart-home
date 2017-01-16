var gulp = require('gulp'),
    conf = require('../config').build_templates,
    connect = require('gulp-connect'),
    templateCache = require('gulp-angular-templatecache');

gulp.task('build_templates', function() {
    return gulp.src(conf.source)
        .pipe(templateCache(conf.filename, {
            transformUrl: function(url) {
                return conf.prefix + url.replace(/\.html\.html/, '.html')
            }
        }))
        .pipe(gulp.dest(conf.dest))
        .pipe(connect.reload());
});