const ffi = require('ffi-napi');
const ref = require('ref-napi');

const StructType = require('ref-struct-napi');

const GravityError = StructType({
	message: ref.types.CString,
});

const GravityErrorPtr = ref.refType(GravityError);

// SubscriberOptions
const SubscriberOptions = StructType({
	endpoint: ref.types.CString,
	domain: ref.types.CString,
	workerCount: ref.types.int,
	bufferSize: ref.types.int,
	chunkSize: ref.types.int,
	verbose: ref.types.bool,
});

const SubscriberOptionsPtr = ref.refType(SubscriberOptions);

const SubscriberMessage = StructType({
});

const SubscriberMessagePtr = ref.refType(SubscriberMessage);

const gravity = ffi.Library('../gravity-sdk', {
	'NewSubscriberOptions': [ SubscriberOptionsPtr, [] ],
	'NewSubscriber': [ 'pointer', [ SubscriberOptionsPtr ] ],
	'SubscriberConnect': [ GravityErrorPtr, [ 'pointer', ref.types.CString ] ],
	'SubscriberSetEventHandler': [ 'void', [ 'pointer', 'pointer' ] ],
	'SubscriberSetSnapshotHandler': [ 'void', [ 'pointer', 'pointer' ] ],
	'SubscriberRegister': [ GravityErrorPtr, [ 'pointer', ref.types.CString, ref.types.CString, ref.types.CString, ref.types.CString ] ],
})


// Create a new subscriber
let opts = gravity.NewSubscriberOptions().deref()

console.log(opts.endpoint)
console.log(opts.domain)
console.log(opts.workerCount)
console.log(opts.bufferSize)
console.log(opts.chunkSize)
console.log(opts.verbose)

let subscriber = gravity.NewSubscriber(opts.ref())

let err = gravity.SubscriberConnect(subscriber, '0.0.0.0:32803')
if (!ref.isNull(err)) {
	console.log(err.deref());
	return;
}

// Event handler
let eventHandler = ffi.Callback('void', [ SubscriberMessagePtr ], function(msg) {
	console.log(msg.deref())
});

gravity.SubscriberSetEventHandler(subscriber, eventHandler);

// Snapshot handler
let snapshotHandler = ffi.Callback('void', [ SubscriberMessagePtr ], function(msg) {
	console.log(msg.deref())
});

gravity.SubscriberSetSnapshotHandler(subscriber, snapshotHandler);

// Register
e = gravity.SubscriberRegister(subscriber, 'transmitter', 'testing', 'Node.js-Client', 'Node.js Client');
if (!ref.isNull(e)) {
	console.log(e.deref())
	return;
}
