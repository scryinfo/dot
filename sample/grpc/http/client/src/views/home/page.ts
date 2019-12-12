import {createModule} from 'direct-vuex';
import {moduleActionContext} from "./store";

export interface Page {
    pagesSize: number;
    currentPage: number;
}

const page = createModule({
    state: {
        pagesSize: 10,
        currentPage: 1,
    } as Page,
    actions: {
        SetPage(content, page: number) {
            const { dispatch, commit, getters, state } = pageActionContext(content);
            // Here, 'dispatch', 'commit', 'getters' and 'state' are typed.
            commit.SET_PAGE(page);
        }
    },
    mutations: {
        SET_PAGE: (state, page: number) => {
            if (state.pagesSize > page) {
                state.currentPage = page;
            }
        }
    },
});
export default page;
export const pageActionContext = (context: any) => moduleActionContext(context, page);
