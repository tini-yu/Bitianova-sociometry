export namespace service {
	
	export class ListFile {
	    labels: string[];
	    uuid: string;
	
	    static createFrom(source: any = {}) {
	        return new ListFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.labels = source["labels"];
	        this.uuid = source["uuid"];
	    }
	}
	export class Results {
	    labels: string[];
	    uuid: string;
	    maps: any[];
	
	    static createFrom(source: any = {}) {
	        return new Results(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.labels = source["labels"];
	        this.uuid = source["uuid"];
	        this.maps = source["maps"];
	    }
	}

}

