// @generated by protoc-gen-connect-es v0.13.0 with parameter "target=ts"
// @generated from file api/v1/api.proto (package api.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { ConnectToStreamRequest, ConnectToStreamResponse, CreateEventRequest, CreateEventResponse, GetCommentsRequest, GetCommentsResponse, GetEventsRequest, GetEventsResponse, SendCommentRequest, SendCommentResponse } from "./api_pb.ts";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service api.v1.APIService
 */
export const APIService = {
  typeName: "api.v1.APIService",
  methods: {
    /**
     * @generated from rpc api.v1.APIService.CreateEvent
     */
    createEvent: {
      name: "CreateEvent",
      I: CreateEventRequest,
      O: CreateEventResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc api.v1.APIService.GetEvents
     */
    getEvents: {
      name: "GetEvents",
      I: GetEventsRequest,
      O: GetEventsResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc api.v1.APIService.SendComment
     */
    sendComment: {
      name: "SendComment",
      I: SendCommentRequest,
      O: SendCommentResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc api.v1.APIService.GetComments
     */
    getComments: {
      name: "GetComments",
      I: GetCommentsRequest,
      O: GetCommentsResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc api.v1.APIService.ConnectToStream
     */
    connectToStream: {
      name: "ConnectToStream",
      I: ConnectToStreamRequest,
      O: ConnectToStreamResponse,
      kind: MethodKind.ServerStreaming,
    },
  }
} as const;

