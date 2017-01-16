var gulp = require('gulp'),
    conf = require('../config').build_coffee_js,
    concat = require('gulp-concat'),
    ngconcat = require('gulp-ngconcat'),
    inject = require('gulp-inject'),
    coffee = require('gulp-coffee'),
    uglify = require('gulp-uglify'),
    replace = require('gulp-replace');

var date = new Date();

for (var key in conf) {

    (function (_task, _source, _filename, _dest) {

        var task = 'coffee:' + _task,
            source = _source,
            filename = _filename,
            dest = _dest;

        gulp.task(task, function(done) {
            return gulp.src(source)
                .pipe(coffee({bare: true})
                    .on('error', done))     // Compile coffeescript
                //.pipe(ngconcat(conf.filename))
                .pipe(concat(filename))
                //.pipe(uglify())
                //.pipe(ngClassify())
                .pipe(replace('__CURRENT_TIME__', date))
                .pipe(gulp.dest(dest));

        });

    })(key, conf[key].source, conf[key].filename, conf[key].dest)
}

gulp.task('coffee', ['coffee:public', 'coffee:private']);
