import path from 'path';
import glob from 'glob';
import Vue from 'vue';

export default function beforeAllTests () {
    // Automatically register all components
    const fileComponents = glob.sync(path.join(__dirname, '../components/**/*.vue'));
    for (const file of fileComponents) {
        const name = file.match(/(\w*)\.vue$/)[1];
        Vue.component(name, require(file).default);
    }
}