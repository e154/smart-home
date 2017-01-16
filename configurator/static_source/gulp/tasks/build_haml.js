var gulp = require('gulp'),
    conf = require('../config').build_haml,
    cache = require('gulp-cached'),
    haml = require('gulp-haml');

for (var key in conf) {

    (function (_task, _source, _dest) {

        var task = 'haml:' + _task,
            source = _source,
            dest = _dest;

        gulp.task(task, function() {
            return gulp.src(source, {read: true})
                .pipe(cache('linting'))
                .pipe(haml())
                .pipe(gulp.dest(dest));
            //.pipe(connect.reload());
        });


    })(key, conf[key].source, conf[key].dest)
}

gulp.task('haml', ['haml:public', 'haml:private']);
