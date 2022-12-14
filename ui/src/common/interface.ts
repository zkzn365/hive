export interface FormValue<T = any> {
  value: T;
  isInvalid: boolean;
  errorMsg: string;
  [prop: string]: any;
}

export interface FormDataType {
  [prop: string]: FormValue;
}

export interface Paging {
  page: number;
  page_size?: number;
}

export type ReportType = 'question' | 'answer' | 'comment' | 'user';
export type ReportAction = 'close' | 'flag' | 'review';
export interface ReportParams {
  type: ReportType;
  action: ReportAction;
}

export interface TagBase {
  display_name: string;
  slug_name: string;
  recommend: boolean;
  reserved: boolean;
}

export interface Tag extends TagBase {
  main_tag_slug_name?: string;
  original_text?: string;
  parsed_text?: string;
}

export interface SynonymsTag extends Tag {
  tag_id: string;
  tag?: string;
}

export interface TagInfo extends TagBase {
  tag_id: string;
  original_text: string;
  parsed_text: string;
  follow_count: number;
  question_count: number;
  is_follower: boolean;
  member_actions;
  created_at?;
  updated_at?;
  main_tag_slug_name?: string;
  excerpt?;
}
export interface QuestionParams {
  title: string;
  content: string;
  html: string;
  tags: Tag[];
}

export interface ListResult<T = any> {
  count: number;
  list: T[];
}

export interface AnswerParams {
  content: string;
  html: string;
  question_id: string;
  id: string;
  edit_summary?: string;
}

export interface LoginReqParams {
  e_mail: string;
  /** password */
  pass: string;
  captcha_id?: string;
  captcha_code?: string;
}

export interface RegisterReqParams extends LoginReqParams {
  name: string;
}

export interface ModifyPasswordReq {
  old_pass: string;
  pass: string;
}

/** User  */
export interface ModifyUserReq {
  display_name: string;
  username?: string;
  avatar: any;
  bio: string;
  bio_html?: string;
  location: string;
  website: string;
}

export interface UserInfoBase {
  avatar: any;
  username: string;
  display_name: string;
  rank: number;
  website: string;
  location: string;
  ip_info?: string;
  /** 'forbidden' | 'normal' | 'delete'
   */
  status?: string;
  /** roles */
  is_admin?: boolean;
}

export interface UserInfoRes extends UserInfoBase {
  bio: string;
  bio_html: string;
  create_time?: string;
  /**
   * value = 1 active;
   * value = 2 inactivated
   */
  mail_status: number;
  language: string;
  is_admin: boolean;
  e_mail?: string;
  [prop: string]: any;
}

export type UploadType = 'post' | 'avatar' | 'branding';
export interface UploadReq {
  file: FormData;
}

export interface ImgCodeReq {
  captcha_id?: string;
  captcha_code?: string;
}

export interface ImgCodeRes {
  captcha_id: string;
  captcha_img: string;
  verify: boolean;
}

export interface PasswordResetReq extends ImgCodeReq {
  e_mail: string;
}

export interface CheckImgReq {
  action: 'login' | 'e_mail' | 'find_pass';
}

export interface SetNoticeReq {
  notice_switch: boolean;
}

export interface NotificationStatus {
  inbox: number;
  achievement: number;
  revision: number;
  can_revision: boolean;
}

export interface QuestionDetailRes {
  id: string;
  title: string;
  content: string;
  html: string;
  tags: any[];
  view_count: number;
  unique_view_count?: number;
  answer_count: number;
  favorites_count: number;
  follow_counts: 0;
  accepted_answer_id: string;
  last_answer_id: string;
  create_time: string;
  update_time: string;
  user_info: UserInfoBase;
  answered: boolean;
  collected: boolean;

  [prop: string]: any;
}

export interface AnswersReq extends Paging {
  order?: 'default' | 'updated';
  question_id: string;
}

export interface AnswerItem {
  id: string;
  question_id: string;
  content: string;
  html: string;
  create_time: string;
  update_time: string;
  user_info: UserInfoBase;
  [prop: string]: any;
}

export interface PostAnswerReq {
  content: string;
  html: string;
  question_id: string;
}

export interface PageUser {
  id?;
  displayName;
  userName?;
  avatar_url?;
}

export interface LangsType {
  label: string;
  value: string;
}

/**
 * @description interface for Question
 */
export type QuestionOrderBy =
  | 'newest'
  | 'active'
  | 'frequent'
  | 'score'
  | 'unanswered';

export interface QueryQuestionsReq extends Paging {
  order: QuestionOrderBy;
  tag?: string;
}

export type AdminQuestionStatus = 'available' | 'closed' | 'deleted';

export type AdminContentsFilterBy = 'normal' | 'closed' | 'deleted';

export interface AdminContentsReq extends Paging {
  status: AdminContentsFilterBy;
  query?: string;
}

/**
 * @description interface for Answer
 */
export type AdminAnswerStatus = 'available' | 'deleted';

/**
 * @description interface for Users
 */
export type UserFilterBy =
  | 'all'
  | 'staff'
  | 'inactive'
  | 'suspended'
  | 'deleted';

/**
 * @description interface for Flags
 */
export type FlagStatus = 'pending' | 'completed';
export type FlagType = 'all' | 'question' | 'answer' | 'comment';
export interface AdminFlagsReq extends Paging {
  status: FlagStatus;
  object_type: FlagType;
}

/**
 * @description interface for Admin Settings
 */
export interface AdminSettingsGeneral {
  name: string;
  short_description: string;
  description: string;
  site_url: string;
  contact_email: string;
}

export interface AdminSettingsInterface {
  language: string;
  theme: string;
  time_zone?: string;
}

export interface AdminSettingsSmtp {
  encryption: string;
  from_email: string;
  from_name: string;
  smtp_authentication: boolean;
  smtp_host: string;
  smtp_password?: string;
  smtp_port: number;
  smtp_username?: string;
  test_email_recipient?: string;
}

export interface SiteSettings {
  branding: AdmingSettingBranding;
  general: AdminSettingsGeneral;
  interface: AdminSettingsInterface;
}

export interface AdmingSettingBranding {
  logo: string;
  square_icon: string;
  mobile_logo?: string;
  favicon?: string;
}

export interface AdminSettingsLegal {
  privacy_policy_original_text?: string;
  privacy_policy_parsed_text?: string;
  terms_of_service_original_text?: string;
  terms_of_service_parsed_text?: string;
}

export interface AdminSettingsWrite {
  recommend_tags: string[];
  required_tag: string;
  reserved_tags: string[];
}

/**
 * @description interface for Activity
 */
export interface FollowParams {
  is_cancel: boolean;
  object_id: string;
}

/**
 * @description search request params
 */
export interface SearchParams {
  q: string;
  order: string;
  page: number;
  size?: number;
}

/**
 * @description search response data
 */
export interface SearchResItem {
  object_type: string;
  object: {
    id: string;
    question_id?: string;
    title: string;
    excerpt: string;
    created_at: number;
    user_info: UserInfoBase;
    vote_count: number;
    answer_count: number;
    accepted: boolean;
    tags: TagBase[];
    status?: string;
  };
}
export interface SearchRes extends ListResult<SearchResItem> {
  extra: any;
}

export interface AdminDashboard {
  info: {
    question_count: number;
    answer_count: number;
    comment_count: number;
    vote_count: number;
    user_count: number;
    report_count: number;
    uploading_files: boolean;
    smtp: boolean;
    time_zone: string;
    occupying_storage_space: string;
    app_start_time: number;
    https: boolean;
    version_info: {
      remote_version: string;
      version: string;
    };
  };
}

export interface TimelineReq {
  show_vote: boolean;
  object_id: string;
}

export interface TimelineItem {
  activity_id: number;
  revision_id: number;
  created_at: number;
  activity_type: string;
  username: string;
  user_display_name: string;
  comment: string;
  object_id: string;
  object_type: string;
  cancelled: boolean;
  cancelled_at: any;
}

export interface TimelineObject {
  title: string;
  object_type: string;
  question_id: string;
  answer_id: string;
  main_tag_slug_name?: string;
  display_name?: string;
}

export interface TimelineRes {
  object_info: TimelineObject;
  timeline: TimelineItem[];
}

export interface ReviewItem {
  type: 'question' | 'answer' | 'tag';
  info: {
    object_id: string;
    title: string;
    content: string;
    html: string;
    tags: Tag[];
  };
  unreviewed_info: {
    id: string;
    use_id: string;
    object_id: string;
    title: string;
    status: 0 | 1;
    create_at: number;
    user_info: UserInfoBase;
    reason: string;
    content: Tag | QuestionDetailRes | AnswerItem;
  };
}
export interface ReviewResp {
  count: number;
  list: ReviewItem[];
}

export interface UserRoleItem {
  id: number;
  name: string;
  description: string;
}
export interface MemberActionItem {
  action: string;
  name: string;
  type: string;
}
