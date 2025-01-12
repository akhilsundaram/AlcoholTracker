// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {tracker} from '../models';
import {context} from '../models';

export function AddTrackerEntry(arg1:number,arg2:number,arg3:number,arg4:string,arg5:number,arg6:number):Promise<void>;

export function AddTrackerEntryUpdate(arg1:number,arg2:number,arg3:number,arg4:string,arg5:number,arg6:number,arg7:number):Promise<void>;

export function DeleteDrink(arg1:number,arg2:number,arg3:number,arg4:string,arg5:number):Promise<boolean>;

export function GetAlcoholCategories():Promise<Array<string>>;

export function GetDaysSinceLastDrink():Promise<number>;

export function GetDrinkCount(arg1:number,arg2:number,arg3:number):Promise<number>;

export function GetDrinkTagColor(arg1:number,arg2:number,arg3:number):Promise<number>;

export function GetDrinks(arg1:number,arg2:number,arg3:number):Promise<string>;

export function GetEntriesByDate(arg1:number,arg2:number,arg3:number,arg4:string):Promise<Array<tracker.DayData>>;

export function GetEntriesOnDate(arg1:number,arg2:number,arg3:number):Promise<{[key: string]: Array<tracker.DayData>}>;

export function Greet(arg1:string):Promise<string>;

export function Shutdown(arg1:context.Context):Promise<void>;

export function ValidateFormDate(arg1:number,arg2:number,arg3:number):Promise<boolean>;