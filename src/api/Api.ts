/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface ApiCloseGenerationRequestReq {
  status: "completed" | "rejected";
}

export interface ApiLoginRes {
  access_token?: string;
  expires_in?: number;
  token_type?: string;
}

export interface DsCreateTurbine {
  /** @minLength 10 */
  description: string;
  height: number;
  is_active: boolean;
  power: number;
  /**
   * @minLength 3
   * @maxLength 100
   */
  title: string;
}

export interface DsCreateUser {
  /**
   * @minLength 3
   * @maxLength 30
   */
  login: string;
  /**
   * @minLength 8
   * @maxLength 100
   */
  password: string;
}

export interface DsDraftGenerationRequestsBriefInfo {
  generationRequestId?: number;
  turbinesCount?: number;
}

export interface DsGenerationRequest {
  calculated_generation_sum?: number;
  closed_at?: string;
  closed_by?: string;
  created_at?: string;
  created_by?: string;
  formed_at?: string;
  id?: number;
  period_days?: number;
  status?: string;
  turbine_generation_requests?: DsTurbineGenerationRequest[];
  turbine_generation_requests_count?: number;
}

export interface DsTurbine {
  Image: string;
  description?: string;
  height?: number;
  id?: number;
  is_active?: boolean;
  power?: number;
  title?: string;
}

export interface DsTurbineGenerationRequest {
  avg_velocity: number | null | undefined;
  alpha?: number;
  avgVelocity?: number;
  /** @format int64 */
  calculatedGeneration?: number;
  generationRequestID?: number;
  turbine?: DsTurbine;
  turbineID?: number;
}

export interface DsUpdateGenerationRequest {
  period_days?: number;
}

export interface DsUpdateTurbine {
  /** @minLength 10 */
  description?: string;
  height?: number;
  image?: string;
  isActive?: boolean;
  power?: number;
  /**
   * @minLength 3
   * @maxLength 100
   */
  title?: string;
}

export interface DsUpdateTurbineGenerationRequest {
  alpha?: number;
  avg_velocity?: number;
}

export interface DsUpdateUser {
  /**
   * @minLength 3
   * @maxLength 30
   */
  login?: string;
  /**
   * @minLength 8
   * @maxLength 100
   */
  password?: string;
}

export interface DsUser {
  id?: number;
  is_moderator?: boolean;
  login?: string;
}

import type {
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  HeadersDefaults,
  ResponseType,
} from "axios";
import axios from "axios";

export type QueryParamsType = Record<string | number, any>;

export interface FullRequestParams
  extends Omit<AxiosRequestConfig, "data" | "params" | "url" | "responseType"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseType;
  /** request body */
  body?: unknown;
}

export type RequestParams = Omit<
  FullRequestParams,
  "body" | "method" | "query" | "path"
>;

export interface ApiConfig<SecurityDataType = unknown>
  extends Omit<AxiosRequestConfig, "data" | "cancelToken"> {
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void;
  secure?: boolean;
  format?: ResponseType;
}

export enum ContentType {
  Json = "application/json",
  JsonApi = "application/vnd.api+json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
  Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
  public instance: AxiosInstance;
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private secure?: boolean;
  private format?: ResponseType;

  constructor({
    securityWorker,
    secure,
    format,
    ...axiosConfig
  }: ApiConfig<SecurityDataType> = {}) {
    this.instance = axios.create({
      ...axiosConfig,
      baseURL: axiosConfig.baseURL || "http://localhost:8000",
    });
    this.secure = secure;
    this.format = format;
    this.securityWorker = securityWorker;
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected mergeRequestParams(
    params1: AxiosRequestConfig,
    params2?: AxiosRequestConfig,
  ): AxiosRequestConfig {
    const method = params1.method || (params2 && params2.method);

    return {
      ...this.instance.defaults,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...((method &&
          this.instance.defaults.headers[
            method.toLowerCase() as keyof HeadersDefaults
          ]) ||
          {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected stringifyFormItem(formItem: unknown) {
    if (typeof formItem === "object" && formItem !== null) {
      return JSON.stringify(formItem);
    } else {
      return `${formItem}`;
    }
  }

  protected createFormData(input: Record<string, unknown>): FormData {
    if (input instanceof FormData) {
      return input;
    }
    return Object.keys(input || {}).reduce((formData, key) => {
      const property = input[key];
      const propertyContent: any[] =
        property instanceof Array ? property : [property];

      for (const formItem of propertyContent) {
        const isFileType = formItem instanceof Blob || formItem instanceof File;
        formData.append(
          key,
          isFileType ? formItem : this.stringifyFormItem(formItem),
        );
      }

      return formData;
    }, new FormData());
  }

  public request = async <T = any, _E = any>({
    secure,
    path,
    type,
    query,
    format,
    body,
    ...params
  }: FullRequestParams): Promise<AxiosResponse<T>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const responseFormat = format || this.format || undefined;

    if (
      type === ContentType.FormData &&
      body &&
      body !== null &&
      typeof body === "object"
    ) {
      body = this.createFormData(body as Record<string, unknown>);
    }

    if (
      type === ContentType.Text &&
      body &&
      body !== null &&
      typeof body !== "string"
    ) {
      body = JSON.stringify(body);
    }

    return this.instance.request({
      ...requestParams,
      headers: {
        ...(requestParams.headers || {}),
        ...(type ? { "Content-Type": type } : {}),
      },
      params: query,
      responseType: responseFormat,
      data: body,
      url: path,
    });
  };
}

/**
 * @title Turbines API
 * @version 0.2.0
 * @license MIT License
 * @baseUrl http://localhost:8000
 * @contact Егор Бляблин
 *
 * API для работы с сервисом расчета генерации электроэнергии ветрогенераторами
 */
export class Api<
  SecurityDataType extends unknown,
> extends HttpClient<SecurityDataType> {
  api = {
    /**
     * @description Возвращает список заявок пользователя с фильтрацией по дате и статусу (кроме удаленных и черновиков)
     *
     * @tags Заявки расчета выработки
     * @name GenerationRequestsList
     * @summary Список заявок с фильтрацией
     * @request GET:/api/generation-requests/
     * @secure
     */
    generationRequestsList: (
      formedAtBegin?: string,
      formedAtEnd?: string,
      status?: "sent" | "completed" | "rejected",
      params: RequestParams = {},
    ) => {
      const query: Record<string, any> = {};

      if (formedAtBegin) query.formed_at_begin = formedAtBegin;
      if (formedAtEnd) query.formed_at_end = formedAtEnd;
      if (status) query.status = status;

      return this.request<DsGenerationRequest[], Record<string, string>>({
        path: `/api/generation-requests/`,
        method: "GET",
        secure: true,
        query,
        type: ContentType.Json,
        format: "json",
        ...params,
      })
    },

    /**
     * @description Возвращает информацию о черновике заявки текущего пользователя
     *
     * @tags Заявки расчета выработки
     * @name GenerationRequestsDraftList
     * @summary Получение информации о черновике заявки
     * @request GET:/api/generation-requests/draft/
     * @secure
     */
    generationRequestsDraftList: (params: RequestParams = {}) =>
      this.request<DsDraftGenerationRequestsBriefInfo, Record<string, string>>({
        path: `/api/generation-requests/draft/`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Обновляет поля черновика заявки текущего пользователя
     *
     * @tags Заявки расчета выработки
     * @name GenerationRequestsDraftUpdate
     * @summary Обновление черновика заявки
     * @request PUT:/api/generation-requests/draft/
     * @secure
     */
    generationRequestsDraftUpdate: (
      request: DsUpdateGenerationRequest,
      params: RequestParams = {},
    ) =>
      this.request<DsGenerationRequest, Record<string, string>>({
        path: `/api/generation-requests/draft/`,
        method: "PUT",
        body: request,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Удаляет черновик заявки текущего пользователя
     *
     * @tags Заявки расчета выработки
     * @name GenerationRequestsDraftDelete
     * @summary Удаление черновика заявки
     * @request DELETE:/api/generation-requests/draft/
     * @secure
     */
    generationRequestsDraftDelete: (params: RequestParams = {}) =>
      this.request<void, Record<string, string>>({
        path: `/api/generation-requests/draft/`,
        method: "DELETE",
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description Формирует черновик заявки и отправляет на рассмотрение
     *
     * @tags Заявки расчета выработки
     * @name GenerationRequestsDraftSubmitUpdate
     * @summary Отправка черновика заявки
     * @request PUT:/api/generation-requests/draft/submit/
     * @secure
     */
    generationRequestsDraftSubmitUpdate: (params: RequestParams = {}) =>
      this.request<Record<string, string>, Record<string, string>>({
        path: `/api/generation-requests/draft/submit/`,
        method: "PUT",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Обновляет параметры (avg_velocity, alpha) турбины в черновике
     *
     * @tags Заявки расчета выработки
     * @name GenerationRequestsDraftUpdate2
     * @summary Обновление параметров турбины в черновике
     * @request PUT:/api/generation-requests/draft/{turbineId}/
     * @originalName generationRequestsDraftUpdate
     * @duplicate
     * @secure
     */
    generationRequestsDraftUpdate2: (
      turbineId: number,
      updates: DsUpdateTurbineGenerationRequest,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, string>, Record<string, string>>({
        path: `/api/generation-requests/draft/${turbineId}/`,
        method: "PUT",
        body: updates,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Добавляет турбину в черновик заявки текущего пользователя
     *
     * @tags Заявки расчета выработки
     * @name GenerationRequestsDraftCreate
     * @summary Добавление турбины в черновик
     * @request POST:/api/generation-requests/draft/{turbineId}/
     * @secure
     */
    generationRequestsDraftCreate: (
      turbineId: number,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, string>, Record<string, string>>({
        path: `/api/generation-requests/draft/${turbineId}/`,
        method: "POST",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Удаляет турбину из черновика заявки текущего пользователя
     *
     * @tags Заявки расчета выработки
     * @name GenerationRequestsDraftDelete2
     * @summary Удаление турбины из черновика
     * @request DELETE:/api/generation-requests/draft/{turbineId}/
     * @originalName generationRequestsDraftDelete
     * @duplicate
     * @secure
     */
    generationRequestsDraftDelete2: (
      turbineId: number,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, string>, Record<string, string>>({
        path: `/api/generation-requests/draft/${turbineId}/`,
        method: "DELETE",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Возвращает информацию о конкретной заявке с деталями и связанными турбинами
     *
     * @tags Заявки расчета выработки
     * @name GenerationRequestsList2
     * @summary Получение информации о заявке
     * @request GET:/api/generation-requests/{generationRequestId}/
     * @originalName generationRequestsList
     * @duplicate
     * @secure
     */
    generationRequestsList2: (
      generationRequestId: number,
      params: RequestParams = {},
    ) =>
      this.request<DsGenerationRequest, Record<string, string>>({
        path: `/api/generation-requests/${generationRequestId}/`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Модератор завершает или отклоняет заявку (только для модераторов)
     *
     * @tags Заявки расчета выработки
     * @name GenerationRequestsCloseUpdate
     * @summary Завершение/отклонение заявки модератором
     * @request PUT:/api/generation-requests/{generationRequestId}/close/
     * @secure
     */
    generationRequestsCloseUpdate: (
      generationRequestId: number,
      request: ApiCloseGenerationRequestReq,
      params: RequestParams = {},
    ) =>
      this.request<DsGenerationRequest, Record<string, string>>({
        path: `/api/generation-requests/${generationRequestId}/close/`,
        method: "PUT",
        body: request,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Возвращает список активных турбин с возможностью фильтрации по названию
     *
     * @tags Ветрогенераторы
     * @name TurbinesList
     * @summary Список турбин с фильтрацией
     * @request GET:/api/turbines/
     */
    turbinesList: (
      query?: {
        /** Фильтр по названию турбины */
        turbineTitle?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<DsTurbine[], Record<string, string>>({
        path: `/api/turbines/`,
        method: "GET",
        query: query,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Создает новую турбину без изображения
     *
     * @tags Ветрогенераторы
     * @name TurbinesCreate
     * @summary Создание новой турбины
     * @request POST:/api/turbines/
     * @secure
     */
    turbinesCreate: (turbine: DsCreateTurbine, params: RequestParams = {}) =>
      this.request<DsTurbine, Record<string, string>>({
        path: `/api/turbines/`,
        method: "POST",
        body: turbine,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Возвращает информацию о конкретной турбине по ID
     *
     * @tags Ветрогенераторы
     * @name TurbinesList2
     * @summary Получение информации о турбине
     * @request GET:/api/turbines/{turbineId}/
     * @originalName turbinesList
     * @duplicate
     */
    turbinesList2: (turbineId: number, params: RequestParams = {}) =>
      this.request<DsTurbine, Record<string, string>>({
        path: `/api/turbines/${turbineId}/`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Обновляет информацию о турбине по ID
     *
     * @tags Ветрогенераторы
     * @name TurbinesUpdate
     * @summary Обновление информации о турбине
     * @request PUT:/api/turbines/{turbineId}/
     * @secure
     */
    turbinesUpdate: (
      turbineId: number,
      turbine: DsUpdateTurbine,
      params: RequestParams = {},
    ) =>
      this.request<DsTurbine, Record<string, string>>({
        path: `/api/turbines/${turbineId}/`,
        method: "PUT",
        body: turbine,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Удаляет турбину по ID
     *
     * @tags Ветрогенераторы
     * @name TurbinesDelete
     * @summary Удаление турбины
     * @request DELETE:/api/turbines/{turbineId}/
     * @secure
     */
    turbinesDelete: (turbineId: number, params: RequestParams = {}) =>
      this.request<void, Record<string, string>>({
        path: `/api/turbines/${turbineId}/`,
        method: "DELETE",
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description Добавляет или обновляет изображение для турбины
     *
     * @tags Ветрогенераторы
     * @name TurbinesUploadImageCreate
     * @summary Добавление изображения к турбине
     * @request POST:/api/turbines/{turbineId}/upload-image/
     * @secure
     */
    turbinesUploadImageCreate: (
      turbineId: number,
      data: {
        /** Файл изображения */
        image: File;
      },
      params: RequestParams = {},
    ) =>
      this.request<DsTurbine, Record<string, string>>({
        path: `/api/turbines/${turbineId}/upload-image/`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.FormData,
        format: "json",
        ...params,
      }),

    /**
     * @description Возвращает информацию о текущем аутентифицированном пользователе
     *
     * @tags Пользователи
     * @name UsersList
     * @summary Получение информации о текущем пользователе
     * @request GET:/api/users/
     * @secure
     */
    usersList: (params: RequestParams = {}) =>
      this.request<DsUser, Record<string, string>>({
        path: `/api/users/`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Обновляет информацию о текущем аутентифицированном пользователе
     *
     * @tags Пользователи
     * @name UsersUpdate
     * @summary Обновление информации о текущем пользователе
     * @request PUT:/api/users/
     * @secure
     */
    usersUpdate: (user: DsUpdateUser, params: RequestParams = {}) =>
      this.request<DsUser, Record<string, string>>({
        path: `/api/users/`,
        method: "PUT",
        body: user,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Создает нового пользователя с указанными учетными данными
     *
     * @tags Пользователи
     * @name UsersCreate
     * @summary Регистрация нового пользователя
     * @request POST:/api/users/
     */
    usersCreate: (user: DsCreateUser, params: RequestParams = {}) =>
      this.request<DsUser, Record<string, string>>({
        path: `/api/users/`,
        method: "POST",
        body: user,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Аутентифицирует пользователя и создает сессию
     *
     * @tags Пользователи
     * @name UsersLoginCreate
     * @summary Аутентификация пользователя
     * @request POST:/api/users/login/
     */
    usersLoginCreate: (credentials: DsCreateUser, params: RequestParams = {}) =>
      this.request<ApiLoginRes, void>({
        path: `/api/users/login/`,
        method: "POST",
        body: credentials,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Завершает сессию текущего аутентифицированного пользователя, добавляя JWT в черный список
     *
     * @tags Пользователи
     * @name UsersLogoutCreate
     * @summary Деавторизация пользователя
     * @request POST:/api/users/logout/
     * @secure
     */
    usersLogoutCreate: (params: RequestParams = {}) =>
      this.request<Record<string, string>, any>({
        path: `/api/users/logout/`,
        method: "POST",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
}
