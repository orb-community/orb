function tagComparator(propA, propB) {
    const str_a = JSON.stringify(propA);
    const str_b = JSON.stringify(propB);
    // sort by key
    if (str_a.toLowerCase() < str_b.toLowerCase()) {
        return -1;
    }
    if (str_a.toLowerCase() > str_b.toLowerCase()) {
        return 1;
    }
    // eq
    return 0;
}

export { tagComparator };
