const events = require('events');
const ffi = require('ffi-napi');
const ref = require('ref-napi');
const utils = require('./utils');
const StructType = require('ref-struct-napi');
const ArrayType = require('ref-array-napi')
const nativeModule = require('./native');

// SubscriberInitialLoadOptions
const SubscriberInitialLoadOptions = StructType({
	enabled: ref.types.bool,
	mode: ref.types.CString,
	omittedCount: ref.types.uint64,
});

const SubscriberInitialLoadOptionsPtr = ref.refType(SubscriberInitialLoadOptions);

// SubscriberOptions
const SubscriberOptions = StructType({
	endpoint: ref.types.CString,
	domain: ref.types.CString,
	workerCount: ref.types.int,
	bufferSize: ref.types.int,
	chunkSize: ref.types.int,
	verbose: ref.types.bool,
	initialLoad: SubscriberInitialLoadOptions,
});

const SubscriberOptionsPtr = ref.refType(SubscriberOptions);

const StringArray = ArrayType(ref.types.CString);

// SubscriberMessage
const SubscriberMessage = StructType({
	instance: 'pointer',
	pipelineID: ref.types.uint64,
	eventName: ref.types.CString,
	collection: ref.types.CString,
	payload: ref.types.CString,
});

const SubscriberMessagePtr = ref.refType(SubscriberMessage);

const Uint64Array = ArrayType(ref.types.uint64);

// Register methods
nativeModule.register({
	'NewSubscriberOptions': [ SubscriberOptionsPtr, [] ],
	'NewSubscriber': [ 'pointer', [ SubscriberOptionsPtr ] ],
	'NewSubscriberWithClient': [ 'pointer', [ 'pointer', SubscriberOptionsPtr ] ],
	'SubscriberConnect': [ utils.ErrorPtr, [ 'pointer', ref.types.CString ] ],
	'SubscriberSetEventHandler': [ 'void', [ 'pointer', 'pointer' ] ],
	'SubscriberSetSnapshotHandler': [ 'void', [ 'pointer', 'pointer' ] ],
	'SubscriberRegister': [ utils.ErrorPtr, [ 'pointer', ref.types.CString, ref.types.CString, ref.types.CString, ref.types.CString ] ],
	'SubscriberStart': [ 'void', [ 'pointer' ] ],
	'SubscriberGetPipelineCount': [ ref.types.int, [ 'pointer' ] ],
	'SubscriberAddAllPipelines': [ utils.ErrorPtr, [ 'pointer' ] ],
	'SubscriberSubscribeToPipelines': [ utils.ErrorPtr, [ 'pointer', Uint64Array ] ],
	'SubscriberSubscribeToCollection': [ utils.ErrorPtr, [ 'pointer', ref.types.CString, StringArray, ref.types.int ] ],
	'SubscriberMessageAck': [ 'void', [ SubscriberMessagePtr ] ],
});

class Message {

	constructor(instance) {
		this.instance = instance;

		// Parsing content
		let data = instance.deref();
		this.eventName = data.eventName;
		this.collection = data.collection;
		this.payload = JSON.parse(data.payload);
	}

	ack() {
		nativeModule.getLibrary().SubscriberMessageAck(this.instance);
	}
}

module.exports = class Subscriber extends events.EventEmitter {

	constructor(client, opts) {
		super();

		this.client = client;

		let gravity = nativeModule.getLibrary();

		// Getting default options
		let sOpts = gravity.NewSubscriberOptions().deref();

		if (opts.initialLoad) {
			let initialLoad = Object.assign(sOpts.initialLoad, opts.initialLoad);
			Object.assign(sOpts, opts, {
				initialLoad: initialLoad
			});
		} else {
			Object.assign(sOpts, opts);
		}

		this.options = sOpts;

		// Initializing instance
		this.instance = gravity.NewSubscriberWithClient(this.client.instance, sOpts.ref());
		this.eventHandler = null;
		this.snapshotHandler = null;

		this.init();
	}

	init() {

		let gravity = nativeModule.getLibrary();

		// Event handler
		this.eventHandler = ffi.Callback('void', [ SubscriberMessagePtr ], (msg) => {
			setImmediate(() => {
				let m = new Message(msg);
				this.emit('event', m);
			});
		});

		gravity.SubscriberSetEventHandler(this.instance, this.eventHandler);

		// Snapshot handler
		this.snapshotHandler = ffi.Callback('void', [ SubscriberMessagePtr ], (msg) => {
			setImmediate(() => {
				let m = new Message(msg);
				this.emit('snapshot', m);
			});
		});

		gravity.SubscriberSetSnapshotHandler(this.instance, this.snapshotHandler);
	}

	connect() {
		return new Promise((resolve, reject) => {
			nativeModule.getLibrary().SubscriberConnect(this.instance, '0.0.0.0:32803', (err, res) => {
				if (!ref.isNull(res)) {
					return reject(res.deref());
				}

				resolve();
			});
		});
	}

	register(subscriberType, componentName, subscriberID, subscriberName) {

		return new Promise((resolve, reject) => {

			// Register
			let err = nativeModule.getLibrary().SubscriberRegister.async(this.instance, subscriberType, componentName, subscriberID, subscriberName, (err, res) => {
				if (!ref.isNull(res)) {
					return reject(res.deref());
				}
			});

			resolve();
		});
	}

	start() {
		nativeModule.getLibrary().SubscriberStart(this.instance);
	}

	getPipelineCount() {
		return new Promise((resolve, reject) => {
			nativeModule.getLibrary().SubscriberGetPipelineCount.async(this.instance, (err, res) => {
				resolve(res);
			});
		});
	}

	addAllPipelines() {

		return new Promise((resolve, reject) => {
			nativeModule.getLibrary().SubscriberAddAllPipelines.async(this.instance, (err, res) => {
				if (!ref.isNull(res)) {
					return reject(res.deref());
				}

				resolve();
			});
		});
	}

	addPipelines(pipelines) {

		return new Promise((resolve, reject) => {
			let arr = new Uint64Array(pipelines);
			nativeModule.getLibrary().SubscriberSubscribeToPipelines.async(this.instance, arr, (err, res) => {
				if (!ref.isNull(res)) {
					return reject(res.deref());
				}

				resolve();
			});
		});
	}

	subscribeToCollections(collectionEntries = {}) {

		return new Promise((resolve, reject) => {

			let count = Object.keys(collectionEntries).length;

			// No entries
			if (count == 0)
				return resolve();

			let gravity = nativeModule.getLibrary();

			for (let [ collectionName, tables ] of Object.entries(collectionEntries)) {

				let arr = new StringArray(tables);
				gravity.SubscriberSubscribeToCollection.async(this.instance, collectionName, arr, tables.length, (err, res) => {
					if (!ref.isNull(res)) {
						return reject(res.deref());
					}

					resolve();
				});
			}
		});
	}
};
