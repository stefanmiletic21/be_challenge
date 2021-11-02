API Server:

contains two endpoints

    - /range?startHour=1630886400&endHour=1631059200 -> used for getting range of hours you want to see (as defined in assignment)
    - /lastHour -> used for getting (nearly) realtime metric for current hour. Last hour value is calculated once a minute.

If data for some hour is not present it will return -1 for it. Data is not present if worker still hasn't calculated that hour

Worker:

I took assumption that when the hour is over we will not get any new data about that hour.
In that case it is easier to just pre-calculate expenses for every hour and when request comes to the APIServer we 
just query the calculations for each hour. It would be to processing/time consuming to do calculation on request
(for instance if range of month is requested server would timeout).
I made new table aggregated_expenses (if it isn't existing it will be created when we start system).
Worker has two components:

    - one that is recalculating current hour expenses every minute
    - one that starts from the furthest calculated hour and goes back to the configured limit and calculates expenses 
    for every hour
    
Explanation of approach:

As I explained above I feel that calculating on request would be to heavy for the system, and that it would lead to limited
usefulness (you couldn't for instance use the output to make the graph of expenses for the last year).
Second part of worker that calculates past hours could have been migration script, but I fear that it would take too long
in the real usecase, and if anything failed you would need to restart it. With Worker it will go slowly back and when it finishes it will
shutdown.

Possible enhancements and known issues:

    - tests - lack of them doesnt mean I dont think they are important, I just considered this to be more of a POC than production ready server
    - error handling - right now pretty basic, just for initial debuging and catching problems, some logging with context wrapping error would be next
    - support for multiple workers - to improve time of calculation more workers could be added, and that would require transactions and syncing between them
