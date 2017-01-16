var gulp = require('gulp'),
    conf = require('../config').ace_themes;

gulp.task('ace_themes', function(done) {
    return gulp.src(conf.source)
        .pipe(gulp.dest(conf.dest));
});
