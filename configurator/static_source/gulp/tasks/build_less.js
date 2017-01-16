var gulp = require('gulp'),
    conf = require('../config').build_less,
    less = require('gulp-less'),
    concat = require('gulp-concat'),
    gutil = require('gulp-util');

for (var key in conf) {

    (function (_task, _source, _filename, _dest) {

        var task = 'less:' + _task,
            source = _source,
            filename = _filename,
            dest = _dest;

        gulp.task(task, function(done) {
            return gulp.src(source)
                .pipe(concat(filename))
                .pipe(less())
                .on('error', function(err){
                    gutil.log(err);
                    this.emit('end');
                })
                .pipe(gulp.dest(dest));
        });

    })(key, conf[key].source, conf[key].filename, conf[key].dest)
}

gulp.task('less', ['less:public', 'less:private']);
