var source = "static_source";
var tmp = "tmp";
var assets = "assets";

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
    },
    redactor_theme_less: {
        filename: 'style.css',
        watch: 'static_source/app/constructor/themes/**/*',
        minimal: {
            source: [
                'static_source/app/constructor/themes/minimal/**/*.less'
            ],
            dest: 'static_source/themes/minimal'
        }
    },
    redactor_theme_files: {
        watch: 'static_source/app/constructor/themes/**/*',
        minimal: {
            source: [
                'static_source/app/constructor/themes/minimal/**/*.svg'
            ],
            dest: 'static_source/themes/minimal'
        }
    }
};
