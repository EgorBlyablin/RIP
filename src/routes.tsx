export const ROUTES = {
    HOME: "/",
    TURBINES_LIST: "/turbines",
    LOGIN: "/login",
    REGISTRATION: "/registration",
    PROFILE: "/profile",
    GENERATION_REQUESTS: "/generation-requests"
}

export const ROUTE_LABELS: Record<keyof typeof ROUTES, string> = {
    HOME: "Главная",
    TURBINES_LIST: "Ветрогенераторы",
    LOGIN: "Авторизация",
    REGISTRATION: "Регистрация",
    PROFILE: "Профиль",
    GENERATION_REQUESTS: "Заявки расчета"
}