var gulp = require('gulp'),
    conf = require('../config').copy;

for (var key in conf) {

    (function (_task, _source, _dest) {

        var task = 'copy:' + _task,
            source = _source,
            dest = _dest;

        gulp.task(task, function(done) {
            return gulp.src(source)
                .pipe(gulp.dest(dest));
        });

    })(key, conf[key].source, conf[key].dest)
}

gulp.task('copy', ['copy:ace', 'copy:translate', 'copy:font-awesome', 'copy:bpmn-font', 'copy:images']);