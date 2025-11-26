export namespace config {
	
	export class WindowSettings {
	    width: number;
	    height: number;
	    x: number;
	    y: number;
	    maximised: boolean;
	
	    static createFrom(source: any = {}) {
	        return new WindowSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.width = source["width"];
	        this.height = source["height"];
	        this.x = source["x"];
	        this.y = source["y"];
	        this.maximised = source["maximised"];
	    }
	}

}

export namespace file {
	
	export class Filter {
	    displayName: string;
	    pattern: string;
	
	    static createFrom(source: any = {}) {
	        return new Filter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.displayName = source["displayName"];
	        this.pattern = source["pattern"];
	    }
	}
	export class SaveOptions {
	    title: string;
	    defaultFilename: string;
	    filters: Filter[];
	
	    static createFrom(source: any = {}) {
	        return new SaveOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.defaultFilename = source["defaultFilename"];
	        this.filters = this.convertValues(source["filters"], Filter);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace updater {
	
	export class UpdateInfo {
	    currentVersion: string;
	    latestVersion: string;
	    hasUpdate: boolean;
	    description: string;
	    downloadUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currentVersion = source["currentVersion"];
	        this.latestVersion = source["latestVersion"];
	        this.hasUpdate = source["hasUpdate"];
	        this.description = source["description"];
	        this.downloadUrl = source["downloadUrl"];
	    }
	}

}

