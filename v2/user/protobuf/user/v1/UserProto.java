// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: user/v1/user.proto

package user.v1;

public final class UserProto {
  private UserProto() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_user_v1_GetUserinfoRequest_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_user_v1_GetUserinfoRequest_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n\022user/v1/user.proto\022\007user.v1\032\026user/v1/u" +
      "serinfo.proto\"\"\n\022GetUserinfoRequest\022\014\n\004n" +
      "ame\030\001 \001(\t2F\n\004User\022>\n\014get_userinfo\022\033.user" +
      ".v1.GetUserinfoRequest\032\021.user.v1.Userinf" +
      "oB,\n\007user.v1B\tUserProtoP\001Z\024/go/pkg/user/" +
      "v1;userb\006proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
          user.v1.UserinfoProto.getDescriptor(),
        });
    internal_static_user_v1_GetUserinfoRequest_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_user_v1_GetUserinfoRequest_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_user_v1_GetUserinfoRequest_descriptor,
        new java.lang.String[] { "Name", });
    user.v1.UserinfoProto.getDescriptor();
  }

  // @@protoc_insertion_point(outer_class_scope)
}
