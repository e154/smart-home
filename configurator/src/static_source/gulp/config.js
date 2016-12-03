var source = "static_source";
var tmp = "tmp";

module.exports = {
    build_lib_js: {
        filename: 'lib.min.js',
        paths: {
            bowerDirectory: source + '/bower_components',
            bowerrc: '.bowerrc',
            bowerJson: 'bower.json'
        },
        dest: source + '/js'
    },
    build_coffee_js: {
        filename: 'app.min.js',
        source: [
            source + "/app/app.coffee",
            source + "/app/**/*.coffee"
        ],
        watch: source + "/app/**/*.coffee",
        dest: source + '/js'
    },
    build_lib_css: {
        filename: 'lib.min.css',
        paths: {
            bowerDirectory: source + '/bower_components',
            bowerrc: '.bowerrc',
            bowerJson: 'bower.json'
        },
        dest: source + '/css'
    },
    build_less: {
        filename: 'app.min.css',
        source: [
            source + '/less/bootstrap.less',
            source + '/less/bootstrap-theme.less',
            source + '/less/app.less'
        ],
        dest: source + '/css',
        watch: source + '/less/**/*.less'
    },
    webserver: {
        root: [ source + '/templates', source],
        enabled: false,
        livereload: true,
        port: 8088,
        fallback: 'static_source/index.html',
        watch: [
            'static_source/js/*.js',
            'static_source/css/*.css'
        ]
    },
    build_haml: {
        source: 'static_source/app/**/*.haml',
        tmp: 'static_source/tmp/templates',
        watch: [
            "static_source/app/**/*.haml"
        ]
    },
    build_templates: {
        filename: "templates.min.js",
        prefix: '/',
        source: 'static_source/tmp/templates/**/*.html',
        watch: [
            'static_source/tmp/templates/**/*'
        ],
        dest: 'static_source/js'
    },
    ace_themes: {
        source: source + '/bower_components/ace-builds/**/*',
        dest: source + '/js/ace-builds'
    }
};
