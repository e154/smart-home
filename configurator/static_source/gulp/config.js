var source = "./",
    tmp = "tmp",
    public_dest = "../build/public",
    private_dest = "../build/private",
    public_src = "public",
    private_src = "private",
    public_bower = public_src + "/bower_components",
    private_bower = private_src + "/bower_components";

module.exports = {
    build_lib_js: {
        public: {
            filename: 'lib.min.js',
            paths: {
                bowerDirectory: public_bower,
                bowerrc: public_src + '/.bowerrc',
                bowerJson: public_src + '/bower.json'
            },
            dest: public_dest + '/js'
        },
        private: {
            filename: 'lib.min.js',
            paths: {
                bowerDirectory: private_bower,
                bowerrc: private_src + '/.bowerrc',
                bowerJson: private_src + '/bower.json'
            },
            dest: private_dest + '/js'
        }
    },
    build_lib_css: {
        public: {
            filename: 'lib.min.css',
            paths: {
                bowerDirectory: public_bower,
                bowerrc: public_src + '/.bowerrc',
                bowerJson: public_src + '/bower.json'
            },
            dest: public_dest + '/css'
        },
        private: {
            filename: 'lib.min.css',
            paths: {
                bowerDirectory: private_bower,
                bowerrc: private_src + '/.bowerrc',
                bowerJson: private_src + '/bower.json'
            },
            dest: private_dest + '/css'
        }
    },
    build_coffee_js: {
        public: {
            filename: 'app.min.js',
            source: [
                public_src + "/app/app.coffee",
                public_src + "/app/**/*.coffee"
            ],
            watch: public_src + "/app/**/*.coffee",
            dest: public_dest + '/js'
        },
        private: {
            filename: 'app.min.js',
            source: [
                private_src + "/app/app.coffee",
                private_src + "/app/**/*.coffee"
            ],
            watch: private_src + "/app/**/*.coffee",
            dest: private_dest + '/js'
        }
    },
    build_less: {
        public: {
            filename: 'app.min.css',
            source: [
                public_src + '/less/bootstrap.less',
                public_src + '/less/bootstrap-theme.less',
                public_src + '/less/app.less'
            ],
            dest: public_dest + '/css',
            watch: public_src + '/less/**/*.less'
        },
        private: {
            filename: 'app.min.css',
            source: [
                private_src + '/less/bootstrap.less',
                private_src + '/less/bootstrap-theme.less',
                private_src + '/less/app.less'
            ],
            dest: private_dest + '/css',
            watch: private_src + '/less/**/*.less'
        }

    },
    build_haml: {
        public: {
            source: public_src + '/app/**/*.haml',
            dest: public_src + '/tmp/templates',
            watch: [
                public_src + '/app/**/*.haml'
            ]
        },
        private: {
            source: private_src + '/app/**/*.haml',
            dest: private_src + '/tmp/templates',
            watch: [
                private_src + '/app/**/*.haml'
            ]
        }
    },
    build_templates: {
       public: {
           filename: "templates.min.js",
           prefix: '/',
           source: public_src + '/tmp/templates/**/*.html',
           watch: [
               public_src + '/tmp/templates/**/*'
           ],
           dest: public_dest + '/js'
        },
        private: {
            filename: "templates.min.js",
            prefix: '/',
            source: private_src + '/tmp/templates/**/*.html',
            watch: [
                private_src + '/tmp/templates/**/*'
            ],
            dest: private_dest + '/js'
        }
    },
    redactor_theme_less: {
        filename: 'style.css',
        watch: private_src + '/app/constructor/themes/**/*',
        minimal: {
            source: [
                private_src + '/app/constructor/themes/minimal/**/*.less'
            ],
            dest: private_dest + '/themes/minimal'
        }
    },
    redactor_theme_files: {
        watch: private_src + '/app/constructor/themes/**/*',
        minimal: {
            source: [
                private_src + '/app/constructor/themes/minimal/**/*.svg'
            ],
            dest: private_dest + '/themes/minimal'
        }
    },
    copy: {
        ace: {
            source: private_bower + '/ace-builds/**/*',
            dest: private_dest + '/js/ace-builds'
        },
        translate: {
            source: private_src + '/translates/**/*',
            dest: private_dest + '/translates',
            watch: [
                private_src + '/translates/**/*'
            ]
        },
        'font-awesome': {
            source: private_bower + '/font-awesome/fonts/**/*',
            dest: private_dest + '/fonts'
        },
        'bpmn-font': {
            source: private_bower + '/bpmn-font/dist/font/**/*',
            dest: private_dest + '/font'
        },
        images: {
            source: public_src + '/images/**/*',
            dest: public_dest + '/images'
        }
    }

};
