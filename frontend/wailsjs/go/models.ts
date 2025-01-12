export namespace tracker {
	
	export class DayData {
	    alcohol: string;
	    quantity: number;
	    cost: number;
	    timestamp: number;
	
	    static createFrom(source: any = {}) {
	        return new DayData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.alcohol = source["alcohol"];
	        this.quantity = source["quantity"];
	        this.cost = source["cost"];
	        this.timestamp = source["timestamp"];
	    }
	}

}

