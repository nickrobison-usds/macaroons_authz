//
// Created by Nicholas Robison on 2019-01-29.
// We need to use a custom allocator, in order to currently handle strings.
// I'm not entirely sure how/why this works, but I found this solution:
// https://stackoverflow.com/questions/50215886/boost-interprocess-managed-shared-memory-errors
//

#ifndef CMSAUTHZCLI_SHAREDMEMORY_HPP
#define CMSAUTHZCLI_SHAREDMEMORY_HPP

#include <boost/interprocess/sync/interprocess_semaphore.hpp>

template<typename Alloc>
struct SharedMemory {

    using allocator_type = typename Alloc::template rebind<char>::other;

public:

    explicit SharedMemory(size_t aSize, Alloc alloc = {}) : mutex(0), done(0), running(0), mData(aSize, alloc) {
        // Not used
    }

    std::vector<char, Alloc> mData;
    boost::interprocess::interprocess_semaphore mutex, done, running;
};

#endif //CMSAUTHZCLI_SHAREDMEMORY_HPP
