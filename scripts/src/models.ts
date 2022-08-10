interface Condition {
    fact : "distance" | "weight"
    operator : "lessThan" | "greaterThanOrEqual" | "lessThanOrEqual"
    value : number
  }
  interface Offer {
    code  : string
    discount : number
    conditions: Condition[]
  }
  
  export type Offers = Offer[]
  