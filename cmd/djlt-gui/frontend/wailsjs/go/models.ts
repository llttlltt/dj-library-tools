export namespace config {
	
	export class Endpoint {
	    source_id: string;
	    resource: string;
	    query?: string;
	
	    static createFrom(source: any = {}) {
	        return new Endpoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.source_id = source["source_id"];
	        this.resource = source["resource"];
	        this.query = source["query"];
	    }
	}
	export class Source {
	    id: string;
	    name: string;
	    provider: string;
	    config: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new Source(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.provider = source["provider"];
	        this.config = source["config"];
	    }
	}
	export class Step {
	    id: string;
	    kind: string;
	    source: Endpoint;
	    targets: Endpoint[];
	    after?: string[];
	    options?: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new Step(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.kind = source["kind"];
	        this.source = this.convertValues(source["source"], Endpoint);
	        this.targets = this.convertValues(source["targets"], Endpoint);
	        this.after = source["after"];
	        this.options = source["options"];
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
	export class UpdateConfig {
	    last_check_at: string;
	    check_interval_hour: number;
	
	    static createFrom(source: any = {}) {
	        return new UpdateConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.last_check_at = source["last_check_at"];
	        this.check_interval_hour = source["check_interval_hour"];
	    }
	}
	export class Workflow {
	    id: string;
	    name: string;
	    steps: Step[];
	
	    static createFrom(source: any = {}) {
	        return new Workflow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.steps = this.convertValues(source["steps"], Step);
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

export namespace gui {
	
	export class GroupRow {
	    id: string;
	    name: string;
	    kind: string;
	    parent: string;
	    items: number;
	
	    static createFrom(source: any = {}) {
	        return new GroupRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.kind = source["kind"];
	        this.parent = source["parent"];
	        this.items = source["items"];
	    }
	}
	export class PlexAuthChallenge {
	    url: string;
	    pin_id: number;
	
	    static createFrom(source: any = {}) {
	        return new PlexAuthChallenge(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.pin_id = source["pin_id"];
	    }
	}
	export class TrackRow {
	    id: string;
	    title: string;
	    artist: string;
	    bpm: string;
	
	    static createFrom(source: any = {}) {
	        return new TrackRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.artist = source["artist"];
	        this.bpm = source["bpm"];
	    }
	}
	export class QueryResult {
	    kind: string;
	    tracks: TrackRow[];
	    groups: GroupRow[];
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new QueryResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.kind = source["kind"];
	        this.tracks = this.convertValues(source["tracks"], TrackRow);
	        this.groups = this.convertValues(source["groups"], GroupRow);
	        this.count = source["count"];
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
	export class StepDiff {
	    step_id: string;
	    target_name: string;
	    current: TrackRow[];
	    added: TrackRow[];
	    removed: TrackRow[];
	    unchanged: TrackRow[];
	
	    static createFrom(source: any = {}) {
	        return new StepDiff(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.step_id = source["step_id"];
	        this.target_name = source["target_name"];
	        this.current = this.convertValues(source["current"], TrackRow);
	        this.added = this.convertValues(source["added"], TrackRow);
	        this.removed = this.convertValues(source["removed"], TrackRow);
	        this.unchanged = this.convertValues(source["unchanged"], TrackRow);
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

export namespace update {
	
	export class UpdateInfo {
	    available: boolean;
	    version: string;
	    current: string;
	    release_notes: string;
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.available = source["available"];
	        this.version = source["version"];
	        this.current = source["current"];
	        this.release_notes = source["release_notes"];
	        this.url = source["url"];
	    }
	}

}

export namespace workflow {
	
	export class StepResult {
	    step_id: string;
	    status: string;
	    previews?: string[];
	    successes?: string[];
	    warnings?: string[];
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new StepResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.step_id = source["step_id"];
	        this.status = source["status"];
	        this.previews = source["previews"];
	        this.successes = source["successes"];
	        this.warnings = source["warnings"];
	        this.error = source["error"];
	    }
	}
	export class WorkflowResult {
	    workflow_id: string;
	    steps: StepResult[];
	
	    static createFrom(source: any = {}) {
	        return new WorkflowResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.workflow_id = source["workflow_id"];
	        this.steps = this.convertValues(source["steps"], StepResult);
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

