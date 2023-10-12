import { Injectable } from '@angular/core';
import * as YAML from 'yaml';

@Injectable({
    providedIn: 'root',
})

export class CodeEditorService {

    constructor() {

    }

    isYaml(str: string) {
        try {
            YAML.parse(str);
            return true;
        } catch {
            return false;
        }
    }
    isJson(str: string) {
        try {
            JSON.parse(str);
            return true;
        } catch {
            return false;
        }
    }

    checkEmpty (object) {
        for (const key in object) {
            if (object[key] === '' || typeof object[key] === 'undefined' || object[key] === null) {
                return true;
            }
        }
        return false;
    }
}
